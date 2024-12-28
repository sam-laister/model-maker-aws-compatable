package router

import (
	"github.com/Soup666/diss-api/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authController *controller.AuthController,
	taskController *controller.TaskController,
	uploadController *controller.UploadController,
) *gin.Engine {
	// Create a new Gin router
	r := gin.Default()

	// Set up the authentication routes
	r.GET("/login", authController.Login)

	// Tasks
	r.GET("/tasks", taskController.GetTasks)
	r.POST("/tasks", taskController.CreateTask)
	r.POST("/tasks/:taskID/upload", taskController.UploadFileToTask)

	// Uploads
	r.POST("/uploads", uploadController.UploadFile)
	r.GET("/uploads/:filename", uploadController.GetFile)

	return r
}
