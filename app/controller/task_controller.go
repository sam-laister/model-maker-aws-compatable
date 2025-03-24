package controller

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Soup666/diss-api/database"
	"github.com/Soup666/diss-api/model"
	services "github.com/Soup666/diss-api/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	TaskService    services.TaskService
	AppFileService services.AppFileService
	VisionService  services.VisionService
}

func NewTaskController(taskService services.TaskService, appFileService services.AppFileService, visionService services.VisionService) *TaskController {
	return &TaskController{TaskService: taskService, AppFileService: appFileService, VisionService: visionService}
}

func (c *TaskController) GetTasks(ctx *gin.Context) {

	user := ctx.MustGet("user")
	userId := user.(*model.User).Id

	tasks, err := c.TaskService.GetTasks(userId)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	// I dont like this being in a controller, but time constraints
	if err := database.DB.Preload("Mesh", "file_type = ?", "mesh").Find(&tasks).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	ctx.JSON(200, gin.H{"tasks": tasks})
}

func (c *TaskController) GetTask(ctx *gin.Context) {

	// Get the Task ID from the route
	taskIDParam := ctx.Param("taskID")
	taskID, err := strconv.Atoi(taskIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := c.TaskService.GetTask(uint(taskID))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	// Load relations
	c.TaskService.FullyLoadTask(task)

	ctx.JSON(200, gin.H{"task": task})
}

// CreateTask handles task creation
// @Summary Create a new task
// @Description Creates a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Task data"
// @Security BearerAuth
// @Success 201 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (c *TaskController) CreateTask(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	task := &model.Task{
		Title:       "",
		Description: "", // Overriden by ai-description
		UserId:      user.Id,
		Completed:   false,
		Status:      "INITIAL",
	}

	createdTask, err := c.TaskService.CreateTask(task)

	if err != nil {
		log.Printf("Error creating task: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"task": createdTask})
}

// UploadFileToTask handles file uploads for a task
// @Summary Upload files to a task
// @Description Uploads files to a task
// @Tags tasks
// @Accept json
// @Produce json
// Security BearerAuth
// @Param taskID path string true "Task ID"
// @Param files formData file true "Files to upload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{taskID}/upload [post]
func (c *TaskController) UploadFileToTask(ctx *gin.Context) {

	// Get the Task ID from the route
	taskIdParam := ctx.Param("taskID")
	taskId, err := strconv.Atoi(taskIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Check if the Task exists
	var task model.Task
	if err := database.DB.First(&task, taskId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Retrieve files from the request
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	// Define the upload folder
	folderPath := fmt.Sprintf("uploads/%d", taskId)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	var uploadedImages []model.AppFile
	var wg sync.WaitGroup
	var mu sync.Mutex
	var hasError bool

	// Use transactions for better consistency
	tx := database.DB.Begin()

	for index, file := range files {
		wg.Add(1)
		go func(index int, file *multipart.FileHeader) {
			defer wg.Done()

			// Validate file extension
			fileExt := strings.ToLower(filepath.Ext(file.Filename))
			if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
				hasError = true
				return
			}

			// Generate a unique filename
			filename := fmt.Sprintf("%d-%d%s", taskId, index, fileExt)
			savePath := filepath.Join(folderPath, filename)

			// Save the file
			if err := ctx.SaveUploadedFile(file, savePath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file %s", file.Filename)})
				hasError = true
				return
			}

			// Save metadata to DB
			image := model.AppFile{
				Filename: filename,
				Url:      fmt.Sprintf("/uploads/%d/%s", taskId, filename),
				TaskId:   uint(taskId),
				FileType: "upload",
			}

			mu.Lock()
			if err := tx.Create(&image).Error; err != nil {
				hasError = true
			} else {
				uploadedImages = append(uploadedImages, image)
			}
			mu.Unlock()
		}(index, file)
	}

	wg.Wait()

	if hasError {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload some files"})
		return
	}

	tx.Commit()

	go func() {
		// Generate caption
		result, err := c.VisionService.AnalyseImage(fmt.Sprintf("./uploads/%d/%s", taskId, uploadedImages[0].Filename), "")

		if err != nil {
			log.Printf("Unable to analyze the image: %v", err)
			return
		}

		if err := c.TaskService.UpdateMeta(&task, "ai-description", result); err != nil {
			log.Printf("Failed to update task metadata: %v", err)
		}
	}()

	go func() {
		// Generate caption
		result, err := c.VisionService.AnalyseImage(fmt.Sprintf("./uploads/%d/%s", taskId, uploadedImages[0].Filename), "categorize the model in this image, use one word only")

		if err != nil {
			log.Printf("Unable to analyze the image: %v", err)
			return
		}

		if err := c.TaskService.UpdateMeta(&task, "ai-title", result); err != nil {
			log.Printf("Failed to update task metadata: %v", err)
		}
	}()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully",
		"images":  uploadedImages,
	})
}

// StartProcess handles the process of starting the photogrammetry process
// @Summary Upload files to a task
// @Description Uploads files to a task
// @Tags tasks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param request body CreateTaskRequest true "Task data"
// @Param taskID path string true "Task ID"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{taskID}/start [post]
func (c *TaskController) StartProcess(ctx *gin.Context) {

	taskId := ctx.Param("taskID")
	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := c.TaskService.GetTask(uint(taskIdInt))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			log.Printf("Error retrieving task: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	task.Completed = false
	if err := c.TaskService.SaveTask(task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Respond to the client immediately
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Process started."})

	go c.TaskService.RunPhotogrammetryProcess(task)
}
