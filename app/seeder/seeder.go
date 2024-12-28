package main

import (
	"log"
	"os"

	"github.com/Soup666/diss-api/model"
	"github.com/Soup666/diss-api/seeds"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	DB, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	DB.Migrator().DropTable(&model.Task{}, &model.User{}, &model.AppFile{})
	DB.AutoMigrate(&model.Task{}, &model.User{}, &model.AppFile{})

	for _, seed := range seeds.All() {
		if err := seed.Run(DB); err != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, err)
		}
	}
}
