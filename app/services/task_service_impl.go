package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	models "github.com/Soup666/modelmaker/model"
	repositories "github.com/Soup666/modelmaker/repository"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"gorm.io/gorm"
)

type TaskServiceImpl struct {
	taskRepo            repositories.TaskRepository
	appFileService      AppFileService
	chatRepository      repositories.ChatRepository
	notificationService NotificationService
	storageService      StorageService
	jobQueue            chan TaskJob
}

func NewTaskService(
	taskRepo repositories.TaskRepository,
	appFileService AppFileService,
	chatRepository repositories.ChatRepository,
	notificationService NotificationService,
	storageService StorageService,
) TaskServiceImpl {
	return TaskServiceImpl{
		taskRepo:            taskRepo,
		appFileService:      appFileService,
		chatRepository:      chatRepository,
		notificationService: notificationService,
		storageService:      storageService,
		jobQueue:            make(chan TaskJob, 100),
	}
}

func (s *TaskServiceImpl) CreateTask(task *models.Task) error {
	err := s.taskRepo.CreateTask(task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskServiceImpl) GetTask(taskID uint) (*models.Task, error) {
	task, err := s.taskRepo.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskServiceImpl) GetUnarchivedTasks(userID uint) ([]*models.Task, error) {

	tasks, err := s.taskRepo.GetUnarchivedTasks(userID)

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskServiceImpl) GetArchivedTasks(userID uint) ([]*models.Task, error) {

	tasks, err := s.taskRepo.GetArchivedTasks(userID)

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskServiceImpl) UpdateTask(task *models.Task) error {

	err := s.taskRepo.SaveTask(task)

	if err != nil {
		return err
	}
	return nil
}

func (s *TaskServiceImpl) UpdateMeta(task *models.Task, key string, value interface{}) error {
	if task.Metadata == nil {
		task.Metadata = make(map[string]interface{})
	}
	task.Metadata[key] = value
	err := s.UpdateTask(task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskServiceImpl) ArchiveTask(taskID uint) (*models.Task, error) {

	task, err := s.taskRepo.ArchiveTask(taskID)

	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskServiceImpl) UnarchiveTask(taskID uint) (*models.Task, error) {

	task, err := s.taskRepo.UnarchiveTask(taskID)

	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskServiceImpl) SaveTask(task *models.Task) error {
	err := s.taskRepo.SaveTask(task)

	if err != nil {
		return err
	}
	return nil
}

func (s *TaskServiceImpl) FailTask(task *models.Task, message string) error {
	task.Status = models.FAILED
	if err := s.UpdateTask(task); err != nil {
		return err
	}

	if err := s.AddLog(task.ID, message); err != nil {
		log.Printf("Failed to add log: %v\n", err)
	}

	log.Printf("Task %d failed: %s\n", task.ID, message)

	s.notificationService.SendMessage(&models.Notification{
		UserID:  task.UserId,
		Message: "Task failed",
		Title:   task.Title,
	})

	return nil
}

func (s *TaskServiceImpl) GetTaskFiles(taskID uint, fileType string) ([]models.AppFile, error) {
	files, err := s.appFileService.GetTaskFiles(taskID, fileType)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *TaskServiceImpl) GetTaskFile(taskID uint, fileType string) (*models.AppFile, error) {
	file, err := s.appFileService.GetTaskFile(taskID, fileType)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *TaskServiceImpl) FullyLoadTask(task *models.Task) (*models.Task, error) {
	files, err := s.GetTaskFiles(task.ID, "upload")
	if err != nil {
		return nil, err
	}
	task.Images = files

	mesh, err := s.GetTaskFile(task.ID, "mesh")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			task.Mesh = nil
		} else {
			return nil, err
		}
	} else {
		task.Mesh = mesh
	}

	return task, nil
}

func (s *TaskServiceImpl) SendMessage(taskID uint, message string, sender string) (*models.ChatMessage, error) {
	chatMessage := &models.ChatMessage{
		Message: message,
		TaskId:  taskID,
		Sender:  sender,
	}

	err := s.chatRepository.CreateChat(chatMessage)
	if err != nil {
		return chatMessage, err
	}

	return chatMessage, nil
}

func (s *TaskServiceImpl) AddLog(taskID uint, log string) error {

	err := s.taskRepo.AddLog(taskID, log)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskServiceImpl) EnqueueJob(job TaskJob) bool {

	select {
	case s.jobQueue <- job:
		return true
	default:
		return false
	}
}

func (s *TaskServiceImpl) processTask(job TaskJob) {
	ctx := context.Background()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config: %v", err))
	}

	// Create ECS client
	client := ecs.NewFromConfig(cfg)

	bucketInput := fmt.Sprintf("uploads/%d/", job.TaskID)
	if os.Getenv("APP_ENV") == "dev" {
		bucketInput = "development/" + bucketInput
	}

	// Build input for RunTask
	input := &ecs.RunTaskInput{
		Cluster:        aws.String(os.Getenv("AWS_ECS_CLUSTER")),
		LaunchType:     types.LaunchTypeFargate,
		TaskDefinition: aws.String(os.Getenv("AWS_TASK_DEFINITION")),
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        []string{os.Getenv("AWS_SUBNET_ID")},
				SecurityGroups: []string{os.Getenv("AWS_SECURITY_GROUP_ID")},
				AssignPublicIp: types.AssignPublicIpEnabled,
			},
		},
		Tags: []types.Tag{
			{
				Key:   aws.String("task-id"),
				Value: aws.String(fmt.Sprintf("%d", job.TaskID)),
			},
			{
				Key:   aws.String("user-id"),
				Value: aws.String(fmt.Sprintf("%d", job.UserID)),
			},
		},
		Overrides: &types.TaskOverride{
			ContainerOverrides: []types.ContainerOverride{
				{
					Name: aws.String(os.Getenv("AWS_CONTAINER_NAME")),
					Environment: []types.KeyValuePair{
						{Name: aws.String("KATAPULT_ACCESS_KEY"), Value: aws.String(os.Getenv("KATAPULT_ACCESS_KEY"))},
						{Name: aws.String("KATAPULT_BUCKET_NAME"), Value: aws.String(os.Getenv("KATAPULT_BUCKET_NAME"))},
						{Name: aws.String("KATAPULT_ENDPOINT"), Value: aws.String(os.Getenv("KATAPULT_ENDPOINT"))},
						{Name: aws.String("KATAPULT_REGION"), Value: aws.String(os.Getenv("KATAPULT_REGION"))},
						{Name: aws.String("KATAPULT_SECRET_KEY"), Value: aws.String(os.Getenv("KATAPULT_SECRET_KEY"))},
						{Name: aws.String("BUCKET_INPUT"), Value: aws.String(bucketInput)},
						{Name: aws.String("BUCKET_TASK_ID"), Value: aws.String(fmt.Sprintf("%d", job.TaskID))},
					},
				},
			},
		},
	}

	// Run the task
	output, err := client.RunTask(ctx, input)
	if err != nil {
		panic(fmt.Sprintf("failed to run task: %v", err))
	}

	// Print task info
	fmt.Println("Successfully ran task:")
	for _, task := range output.Tasks {
		fmt.Printf("Task ARN: %s\n", *task.TaskArn)
	}
}

func (s *TaskServiceImpl) StartWorker() {
	go func() {
		for job := range s.jobQueue {
			go s.processTask(job)
		}
	}()
}

func (s *TaskServiceImpl) GetJobQueue() chan TaskJob {
	return s.jobQueue
}
