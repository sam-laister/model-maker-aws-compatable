package repositories

import (
	models "github.com/Soup666/diss-api/model"
)

type TaskRepository interface {
	GetTasksByUser(userID uint) ([]models.Task, error)
	GetTaskByID(taskID uint) (*models.Task, error)
	CreateTask(task *models.Task) error
	SaveTask(task *models.Task) error
	ArchiveTask(task *models.Task) error
}
