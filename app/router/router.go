package router

import (
	"github.com/Soup666/diss-api/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authController *controller.AuthController,
	taskController *controller.TaskController,
	uploadController *controller.UploadController,
	objectController *controller.ObjectController,
) *gin.Engine {
	// Create a new Gin router
	r := gin.Default()

	r.Use(CORSMiddleware())

	// Set up the authentication routes
	r.GET("/login", authController.Login)

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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
