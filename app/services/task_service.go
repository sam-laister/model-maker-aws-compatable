package services

import "github.com/Soup666/diss-api/model"

type TaskService interface {
	CreateTask(task *model.Task) (*model.Task, error)
	GetTask(taskID string) (*model.Task, error)
	GetTasks(userID string) ([]model.Task, error)
	UpdateTask(task *model.Task) (*model.Task, error)
	DeleteTask(taskID string) error
	SaveTask(task *model.Task) error
}
