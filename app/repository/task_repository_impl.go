package repositories

import (
	"github.com/Soup666/diss-api/database"
	models "github.com/Soup666/diss-api/model"
	"gorm.io/gorm"
)

type TaskRepositoryImpl struct {
	DB *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &TaskRepositoryImpl{DB: db}
}

func (repo *TaskRepositoryImpl) GetTasksByUser(userID uint) ([]models.Task, error) {
	// Fetch tasks related to the user
	var tasks []models.Task
	if err := database.DB.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (repo *TaskRepositoryImpl) GetTaskByID(taskID uint) (*models.Task, error) {
	var task models.Task
	if err := database.DB.Where("id = ?", taskID).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// CreateTask creates a new task in the database
func (repo *TaskRepositoryImpl) CreateTask(task *models.Task) error {
	if err := database.DB.Create(task).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TaskRepositoryImpl) SaveTask(task *models.Task) error {
	if err := database.DB.Save(task).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TaskRepositoryImpl) ArchiveTask(task *models.Task) error {
	if err := database.DB.Delete(task).Error; err != nil {
		return err
	}
	return nil
}
