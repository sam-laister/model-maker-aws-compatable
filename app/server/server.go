package main

import (
	"context"
	"log"
	"os"

	"github.com/Soup666/diss-api/auth"
	"github.com/Soup666/diss-api/db"
	"github.com/Soup666/diss-api/server/controller"
	"github.com/Soup666/diss-api/server/router"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Set up the database connection
	log.Println("Connecting to database...")
	log.Println(os.Getenv("DATABASE_URL"))
	dbConn, err := db.ConnectToDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	defer dbConn.Close()

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

	// Run database migrations
	dbConn.AutoMigrate(&auth.User{})

	// Set up the authentication service
	authService := &auth.AuthService{
		DB:       dbConn,
		FireAuth: authClient,
	}
	authController := controller.NewAuthController(authService)

	// Set up the HTTP router
	r := router.NewRouter(authController)

	// Start the server
	r.Run(":3333")
}
