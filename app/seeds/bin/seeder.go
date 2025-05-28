package main

import (
	"log"

	"github.com/Soup666/diss-api/database"
	"github.com/Soup666/diss-api/model"
	"github.com/Soup666/diss-api/seeds/seeds"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	err := database.ConnectDatabase()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	log.Default().Println("CreateTestUser")
	if err := seeds.CreateUser(database.DB, "Seed User", "fZ2FW27YhLe6Va7VUQVwYfPlVVU2"); err != nil {
		log.Fatalf("Creating user failed with error: %s", err)
	}

	log.Default().Println("CreateTask")
	task := &model.Task{
		Title:       "Seed Task",
		Description: "This is a seed task",
		Completed:   true,
		UserId:      1,
		Images:      []model.AppFile{},
		Status:      "SUCCESS",
		Metadata:    map[string]interface{}{},
		Mesh:        nil,
	}

	if err := database.DB.Create(task).Error; err != nil {
		log.Fatalf("Creating task failed with error: %s", err)
	}

	log.Default().Println("CreateDummyFiles")
	files, err := seeds.CreateDummyFiles(database.DB, task.ID)
	if err != nil {
		log.Fatalf("Creating dummy files failed with error: %s", err)
	}

	mesh, err := seeds.CreateDummyMesh(database.DB)
	if err != nil {
		log.Fatalf("Creating dummy mesh failed with error: %s", err)
	}

	task.Mesh = mesh
	task.Images = files

	if err := database.DB.Save(task).Error; err != nil {
		log.Fatalf("Saving task with files failed with error: %s", err)
	}

	if err := seeds.CopyRawImages(task.ID); err != nil {
		log.Fatalf("Copying images failed with error: %s", err)
	}

	if err := seeds.CopyRawModel(task.ID); err != nil {
		log.Fatalf("Copying mesh failed with error: %s", err)
	}

	log.Println("Seeding completed successfully")
}
