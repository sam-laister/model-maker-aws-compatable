package router

import (
	"github.com/Soup666/diss-api/controller"
	"github.com/Soup666/diss-api/middleware"
	"github.com/Soup666/diss-api/services"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authController *controller.AuthController,
	taskController *controller.TaskController,
	uploadController *controller.UploadController,
	objectController *controller.ObjectController,
	authService *services.AuthServiceImpl,
) *gin.Engine {

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware(authService))

	// Set up the authentication routes
	r.POST("/verify", authController.Verify)

	// Tasks
	r.GET("/tasks", taskController.GetTasks)
	r.POST("/tasks", taskController.CreateTask)
	r.GET("/tasks/:taskID", taskController.GetTask)
	r.POST("/tasks/:taskID/upload", taskController.UploadFileToTask)
	r.POST("/tasks/:taskID/start", taskController.StartProcess)

	// Uploads
	r.POST("/uploads", uploadController.UploadFile)
	r.GET("/uploads/:taskId/:filename", uploadController.GetFile)

	r.GET("/objects/:taskID/:filename", objectController.GetObject)

	return r
}
