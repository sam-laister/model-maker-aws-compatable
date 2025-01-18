package services

import (
	"github.com/Soup666/diss-api/model"
)

type TaskService interface {
	CreateTask(task *model.Task) (*model.Task, error)
	GetTask(taskID uint) (*model.Task, error)
	GetTasks(userID uint) ([]model.Task, error)
	UpdateTask(task *model.Task) (*model.Task, error)
	DeleteTask(taskID *model.Task) error
	SaveTask(task *model.Task) error
}
