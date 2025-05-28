package repositories_test

import (
	"testing"

	database "github.com/Soup666/diss-api/database"
	"github.com/Soup666/diss-api/model"
	repositories "github.com/Soup666/diss-api/repository"
	"github.com/stretchr/testify/assert"
)

func TestTaskRepository(t *testing.T) {

	err := database.SetupTestDB(t)

	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
		return
	}

	repo := repositories.NewTaskRepository(database.DB)

	// Test Create
	task := &model.Task{
		Title:       "test_task",
		Description: "test_description",
		Status:      "SUCCESS",
	}
	err = repo.CreateTask(task)
	assert.NoError(t, err)
	assert.NotZero(t, task.ID)

	// Test GetTaskByID
	fetchedTask, err := repo.GetTaskByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedTask)
	assert.Equal(t, task.Title, fetchedTask.Title)

	// Test GetTaskByID with non-existent UID
	nonExistentUser, err := repo.GetTaskByID(2)
	assert.Error(t, err)
	assert.Nil(t, nonExistentUser)

	// Test UpdateUser
	task.Title = "test_task2"
	err = repo.SaveTask(task)
	assert.NoError(t, err)

	updatedTask, _ := repo.GetTaskByID(1)
	assert.Equal(t, "test_task2", updatedTask.Title)
}
