package database

import (
	"log"

	models "github.com/Soup666/diss-api/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(connectionString string) {
	var err error
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// DB.Migrator().DropTable(&models.Task{}, &models.User{}, &models.AppFile{})

	// Migrate the schema to ensure tables are created/updated
	DB.AutoMigrate(&models.Task{}, &models.User{}, &models.AppFile{})

}
