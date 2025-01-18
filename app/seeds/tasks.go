package seeds

import (
	"github.com/Soup666/diss-api/model"
	"gorm.io/gorm"
)

func CreateTask(db *gorm.DB, title string, description string, completed bool, userId uint) error {
	return db.Create(&model.Task{
		Title:       title,
		Description: description,
		Completed:   completed,
		UserID:      userId,
		Images:      []model.AppFile{},
		Status:      "INITIAL",
	}).Error
}
