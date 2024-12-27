package repositories

import (
	"github.com/Soup666/diss-api/database"
	models "github.com/Soup666/diss-api/model"
)

// GetUserByAPIKey fetches a user by their API key from the database
func GetTasksByUser(user *models.User) ([]models.Task, error) {
	// Fetch tasks related to the user
	var tasks []models.Task
	if err := database.DB.Where("user_id = ?", user.ID).Find(&tasks).Error; err != nil {
		return nil, nil
	}
	return tasks, nil
}

// CreateTask creates a new task in the database
func CreateTask(task *models.Task) error {
	if err := database.DB.Create(task).Error; err != nil {
		return err
	}
	return nil
}
