package main

import (
	"context"
	"log"
	"os"

	"github.com/Soup666/diss-api/controller"
	db "github.com/Soup666/diss-api/database"
	"github.com/Soup666/diss-api/router"
	"github.com/Soup666/diss-api/services"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	_ "github.com/joho/godotenv/autoload"
)

var authService *services.AuthService

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

	// Set up the authentication service
	authService := &services.AuthService{
		FireAuth: authClient,
	}

	authController := controller.NewAuthController(authService)
	taskController := controller.NewTaskController(authService)
	uploadController := controller.NewUploadController(authService)
	objectController := controller.NewObjectController(authService)

	// Set up the HTTP router
	r := router.NewRouter(authController, taskController, uploadController, objectController)

	// Start the server
	r.Run(":3333")
}
