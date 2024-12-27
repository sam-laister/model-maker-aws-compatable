package controller

import (
	"context"
	"net/http"
	"strings"

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

	ctx.JSON(200, gin.H{"tasks": tasks})
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
