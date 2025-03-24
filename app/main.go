package main

import (
	"context"
	"log"
	"os"

	"github.com/Soup666/diss-api/controller"
	db "github.com/Soup666/diss-api/database"
	repositories "github.com/Soup666/diss-api/repository"
	"github.com/Soup666/diss-api/router"
	"github.com/Soup666/diss-api/services"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Set up the database connection
	log.Println("Connecting to database...")
	log.Println(os.Getenv("DATABASE_URL"))

	db.ConnectDatabase(os.Getenv("DATABASE_URL"))

	// Create a Firebase app instance
	opt := option.WithCredentialsFile("./service-account-key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to create Firebase app: %v", err)
	}

	// Create a Firebase auth client instance
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Firebase auth client: %v", err)
	}

	userRepo := repositories.NewUserRepository(db.DB)
	taskRepo := repositories.NewTaskRepository(db.DB)
	appFileRepo := repositories.NewAppFileRepository(db.DB)

	// Set up the authentication service
	authService := services.NewAuthService(authClient, db.DB, userRepo)
	userService := services.NewUserService(userRepo)
	appFileService := services.NewAppFileServiceFile(appFileRepo)
	taskService := services.NewTaskService(taskRepo, appFileService)
	visionService := services.NewVisionService()

	authController := controller.NewAuthController(authService, userService)
	taskController := controller.NewTaskController(taskService, appFileService, visionService)
	uploadController := controller.NewUploadController()
	objectController := controller.NewObjectController()
	visionController := controller.NewVisionController(visionService, taskRepo, taskService)

	// Set up the HTTP router
	r := router.NewRouter(authController, taskController, uploadController, objectController, visionController, authService)

	// Start the server
	if r.Run(":"+os.Getenv("PORT")) != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}

}
