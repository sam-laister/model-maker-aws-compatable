package main

import (
	"context"
	"log"
	"os"

	"github.com/Soup666/diss-api/auth"
	"github.com/Soup666/diss-api/model"
	"github.com/Soup666/diss-api/server/controller"
	"github.com/Soup666/diss-api/server/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	_ "github.com/joho/godotenv/autoload"
)

var (
	Instance *gorm.DB
)

func dbMigrate() {
	if err := Instance.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
}

func main() {
	// Set up the database connection
	log.Println("Connecting to database...")
	log.Println(os.Getenv("DATABASE_URL"))

	database, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	Instance = database

	if Instance == nil {
		log.Fatal("Database connection is nil")
	}

	Instance.Migrator().DropTable(&model.User{})

	dbMigrate()

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
	authService := &auth.AuthService{
		DB:       Instance,
		FireAuth: authClient,
	}
	authController := controller.NewAuthController(authService)

	// Set up the HTTP router
	r := router.NewRouter(authController)

	// Start the server
	r.Run(":3333")
}
