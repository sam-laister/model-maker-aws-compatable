package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/Soup666/diss-api/database"
	"github.com/Soup666/diss-api/model"
	services "github.com/Soup666/diss-api/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	TaskService    services.TaskService
	AppFileService services.AppFileService
}

func NewTaskController(taskService services.TaskService, appFileService services.AppFileService) *TaskController {
	return &TaskController{TaskService: taskService, AppFileService: appFileService}
}

func (c *TaskController) GetTasks(ctx *gin.Context) {

	user := ctx.MustGet("user")
	userId := user.(*model.User).ID

	tasks, err := c.TaskService.GetTasks(userId)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

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

	if err := database.DB.Preload("Images", "file_type = ?", "upload").First(&task, taskID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Query for the Mesh relation separately
	var mesh *model.AppFile
	if err := database.DB.Where("task_id = ? AND file_type = ?", task.ID, "mesh").First(&mesh).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			task.Mesh = nil
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load Mesh"})
			return
		}
	} else {
		task.Mesh = mesh
	}

	ctx.JSON(200, gin.H{"task": task})
}

// CreateTaskRequest
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CreateTask handles task creation
// @Summary Create a new task
// @Description Creates a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param request body CreateTaskRequest true "Task data"
// @Success 201 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (c *TaskController) CreateTask(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	var req CreateTaskRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		UserID:      user.ID,
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

func (c *TaskController) UploadFileToTask(ctx *gin.Context) {

	// Get the Task ID from the route
	taskIDParam := ctx.Param("taskID")
	taskID, err := strconv.Atoi(taskIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Check if the Task exists
	var task model.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
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

	var uploadedImages []model.AppFile
	folderPath := fmt.Sprintf("uploads/task-%d", taskID)
	os.MkdirAll(folderPath, os.ModePerm)

	for index, file := range files {

		// Generate a unique filename based on the Task ID
		fileType := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("task-%d-%d%s", taskID, index, fileType)
		savePath := filepath.Join(folderPath, filename)

		// Save the file to disk
		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file %s", file.Filename)})
			return
		}

		// Save metadata to the database
		image := model.AppFile{
			Filename: filename,
			Url:      fmt.Sprintf("/uploads/%d/%s", taskID, filename),
			TaskID:   uint(taskID),
			FileType: "upload",
		}
		database.DB.Create(&image)
		uploadedImages = append(uploadedImages, image)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully",
		"images":  uploadedImages,
	})
}

func (c *TaskController) StartProcess(ctx *gin.Context) {

	taskId := ctx.Param("taskID")
	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := c.TaskService.GetTask(uint(taskIdInt))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task not found"})
		return
	}
	task.Completed = false
	if err := c.TaskService.SaveTask(task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Path to the executable
	executablePath := "./cmd/HelloPhotogrammetry"

	var inputPath string = fmt.Sprintf("uploads/task-%s", taskId)
	var buildPath string = fmt.Sprintf("objects/task-%s", taskId)
	var buildFileName string = fmt.Sprintf("baked_mesh_%s.usdz", taskId)

	os.Mkdir(buildPath, os.ModePerm)

	// Create the command
	cmd := exec.Command(executablePath, inputPath, fmt.Sprintf("%s/%s", buildPath, buildFileName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command in a goroutine
	go func() {
		log.Println("Starting process...")
		log.Printf("Command: %v\n", cmd)

		// Start the command
		if err := cmd.Start(); err != nil {
			log.Printf("Failed to start process: %v\n", err)
			return
		}

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			log.Printf("Process finished with error: %v\n", err)
			return
		}

		log.Println("Process completed successfully.")
		task.Completed = true
		if _, err := c.TaskService.UpdateTask(task); err != nil {
			log.Printf("Failed to update task: %v\n", err)
			return
		}

		log.Println("Task updated successfully.")

		var inputPath string = fmt.Sprintf("./objects/task-%d/baked_mesh_%d.usdz", task.ID, task.ID)
		var buildPath string = fmt.Sprintf("./objects/task-%d/task-%d", task.ID, task.ID)

		executablePath := "./cmd/usda_to_glb.sh"
		cmd := exec.Command(executablePath, inputPath, buildPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		go func() {
			log.Println("Starting convertion...")
			log.Printf("Command: %v\n", cmd)

			// Start the command
			if err := cmd.Start(); err != nil {
				log.Printf("Failed to start process: %v\n", err)
				return
			}

			// Wait for the command to finish
			if err := cmd.Wait(); err != nil {
				log.Printf("Process finished with error: %v\n", err)
				return
			}

			log.Println("Process completed successfully.")
			mesh, err := c.AppFileService.Save(&model.AppFile{
				Url:      fmt.Sprintf("/objects/%d/task-%d.glb", task.ID, task.ID),
				Filename: fmt.Sprintf("task-%d.glb", task.ID),
				TaskID:   task.ID,
				FileType: "mesh",
			})

			if err != nil {
				log.Printf("Failed to save mesh: %v\n", err)
				return
			}

			task.Mesh = mesh
			task.Completed = true

			if _, err := c.TaskService.UpdateTask(task); err != nil {
				log.Printf("Failed to update task: %v\n", err)
				return
			}

			log.Println("Task updated successfully.")
		}()
	}()

	// Respond to the client immediately
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Process started."})
}
