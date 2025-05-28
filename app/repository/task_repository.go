package repositories

import (
	models "github.com/Soup666/diss-api/model"
)

type TaskRepository interface {
	GetUnarchivedTasks(userID uint) ([]*models.Task, error)
	GetArchivedTasks(userID uint) ([]*models.Task, error)
	GetTaskByID(taskID uint) (*models.Task, error)
	CreateTask(task *models.Task) error
	SaveTask(task *models.Task) error
	ArchiveTask(task uint) (*models.Task, error)
	UnarchiveTask(task uint) (*models.Task, error)
	AddLog(taskID uint, log string) error
}
