package router

import (
	"github.com/Soup666/diss-api/controller"
	"github.com/Soup666/diss-api/middleware"
	"github.com/Soup666/diss-api/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Soup666/diss-api/docs"
)

func NewRouter(
	authController *controller.AuthController,
	taskController *controller.TaskController,
	uploadController *controller.UploadController,
	objectController *controller.ObjectController,
	authService *services.AuthServiceImpl,
) *gin.Engine {

	r := gin.Default()

	// Global middlewares
	r.Use(middleware.CORSMiddleware())

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Authenticated routes
	authRequired := r.Group("/")
	authRequired.Use(middleware.AuthMiddleware(authService))

	// Authentication routes
	authRequired.POST("/verify", authController.Verify)
	authRequired.PATCH("/verify", authController.Verify)

	// Tasks (protected by AuthMiddleware)
	authRequired.GET("/tasks", taskController.GetTasks)
	authRequired.POST("/tasks", taskController.CreateTask)
	authRequired.GET("/tasks/:taskID", taskController.GetTask)
	authRequired.POST("/tasks/:taskID/upload", taskController.UploadFileToTask)
	authRequired.POST("/tasks/:taskID/start", taskController.StartProcess)

	// Unauthenticated routes
	r.POST("/uploads", uploadController.UploadFile)
	r.GET("/uploads/:taskId/:filename", uploadController.GetFile)
	r.GET("/objects/:taskID/:filename", objectController.GetObject)

	return r
}
