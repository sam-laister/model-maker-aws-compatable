package controller

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Soup666/diss-api/database"
	models "github.com/Soup666/diss-api/model"
	repositories "github.com/Soup666/diss-api/repositories"
	services "github.com/Soup666/diss-api/services"
	"github.com/gin-gonic/gin"
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

	// Preload Images for all tasks
	if err := database.DB.Preload("Images").Find(&tasks).Error; err != nil {
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

	// Use Preload to eagerly load the Images relation
	if err := database.DB.Preload("Images").First(&task, taskID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
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

	task := models.Task{
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
	var task models.Task
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

	var uploadedImages []models.Image
	for _, file := range files {
		// Generate a unique filename based on the Task ID
		filename := fmt.Sprintf("task-%d-%s", taskID, file.Filename)
		savePath := filepath.Join("uploads", filename)

		// Save the file to disk
		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file %s", file.Filename)})
			return
		}

		// Save metadata to the database
		image := models.Image{
			Filename: filename,
			Url:      fmt.Sprintf("/%s", savePath),
			TaskID:   uint(taskID),
		}
		database.DB.Create(&image)
		uploadedImages = append(uploadedImages, image)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully",
		"images":  uploadedImages,
	})
}
