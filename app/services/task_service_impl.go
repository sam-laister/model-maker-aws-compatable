package services

import (
	models "github.com/Soup666/diss-api/model"
	repositories "github.com/Soup666/diss-api/repository"
)

type TaskServiceImpl struct {
	taskRepo repositories.TaskRepository
}

func NewTaskService(taskRepo repositories.TaskRepository) *TaskServiceImpl {
	return &TaskServiceImpl{taskRepo: taskRepo}
}

func (s *TaskServiceImpl) CreateTask(task *models.Task) (*models.Task, error) {
	err := s.taskRepo.CreateTask(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskServiceImpl) GetTask(taskID uint) (*models.Task, error) {
	task, err := s.taskRepo.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskServiceImpl) GetTasks(userID uint) ([]models.Task, error) {

	tasks, err := s.taskRepo.GetTasksByUser(userID)

	if err != nil {
		return nil, err
	}
	return tasks, nil

}

func (s *TaskServiceImpl) UpdateTask(task *models.Task) (*models.Task, error) {

	err := s.taskRepo.SaveTask(task)

	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskServiceImpl) ArchiveTask(taskID uint) error {

	task, err := s.taskRepo.GetTaskByID(taskID)

	if err != nil {
		return err
	}

	err = s.taskRepo.ArchiveTask(task)

	if err != nil {
		return err
	}
	return nil
}

func (s *TaskServiceImpl) SaveTask(task *models.Task) error {
	err := s.taskRepo.SaveTask(task)

	if err != nil {
		return err
	}
	return nil
}

func (s *TaskServiceImpl) DeleteTask(taskID *models.Task) error {
	err := s.taskRepo.ArchiveTask(taskID)

	if err != nil {
		return err
	}
	return nil
}
