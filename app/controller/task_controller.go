package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Soup666/diss-api/database"
	"github.com/Soup666/diss-api/model"
	repositories "github.com/Soup666/diss-api/repositories"
	services "github.com/Soup666/diss-api/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	authService *services.AuthService
}

func NewTaskController(authService *services.AuthService) *TaskController {
	return &TaskController{authService}
}

func (c *TaskController) GetTasks(ctx *gin.Context) {

	// Extract API key from request header
	apiKey := ctx.GetHeader("Authorization")
	if apiKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "API key is missing"})
		return
	}

	// Remove "Bearer " if present
	apiKey = strings.TrimPrefix(apiKey, "Bearer ")

	authToken, err := c.authService.FireAuth.VerifyIDToken(context.Background(), apiKey)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid token"})
		return
	}

	user, err := c.authService.Login(authToken.UID)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tasks, err := repositories.GetTasksByUser(user)
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
	// Extract API key from request header
	apiKey := ctx.GetHeader("Authorization")
	if apiKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "API key is missing"})
		return
	}

	// Remove "Bearer " if present
	apiKey = strings.TrimPrefix(apiKey, "Bearer ")

	_, err := c.authService.FireAuth.VerifyIDToken(context.Background(), apiKey)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid token"})
		return
	}

	// Get the Task ID from the route
	taskIDParam := ctx.Param("taskID")
	taskID, err := strconv.Atoi(taskIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := repositories.GetTaskByID(taskID)
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

func (c *TaskController) CreateTask(ctx *gin.Context) {

	// Extract API key from request header
	apiKey := ctx.GetHeader("Authorization")
	if apiKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "API key is missing"})
		return
	}

	// Remove "Bearer " if present
	apiKey = strings.TrimPrefix(apiKey, "Bearer ")

	authToken, err := c.authService.FireAuth.VerifyIDToken(context.Background(), apiKey)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid token"})
		return
	}

	user, err := c.authService.Login(authToken.UID)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var taskData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&taskData); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if taskData.Title == "" || taskData.Description == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Title and Description are required"})
		return
	}

	task := model.Task{
		Title:       taskData.Title,
		Description: taskData.Description,
		UserID:      user.ID,
		Completed:   false,
	}

	if err := repositories.CreateTask(&task); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"task": task})
}

func (c *TaskController) UploadFileToTask(ctx *gin.Context) {

	// Extract API key from request header
	apiKey := ctx.GetHeader("Authorization")
	if apiKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "API key is missing"})
		return
	}

	// Remove "Bearer " if present
	apiKey = strings.TrimPrefix(apiKey, "Bearer ")

	_, err := c.authService.FireAuth.VerifyIDToken(context.Background(), apiKey)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid token"})
		return
	}

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

	task, err := repositories.GetTaskByID(taskIdInt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task not found"})
		return
	}
	task.Completed = false
	if err := repositories.SaveTask(task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Path to the executable
	executablePath := "./cmd/photogrammetry"

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
		StartConvertion(task)
	}()

	// Respond to the client immediately
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Process started."})
}

func StartConvertion(task *model.Task) {
	task.Completed = true
	if err := repositories.SaveTask(task); err != nil {
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
		mesh, err := repositories.SaveAppFile(&model.AppFile{
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

		if err := repositories.SaveTask(task); err != nil {
			log.Printf("Failed to update task: %v\n", err)
			return
		}

		log.Println("Task updated successfully.")
	}()
}
