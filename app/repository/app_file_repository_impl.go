package repositories

import (
	"github.com/Soup666/diss-api/database"
	"github.com/Soup666/diss-api/model"
	"gorm.io/gorm"
)

type AppFileRepositoryImpl struct {
	DB *gorm.DB
}

func NewAppFileRepository(db *gorm.DB) AppFileRepository {
	return &AppFileRepositoryImpl{DB: db}
}

func (repo *AppFileRepositoryImpl) SaveAppFile(appFile *model.AppFile) (*model.AppFile, error) {
	if err := database.DB.Save(&appFile).Error; err != nil {
		return nil, err
	}
	return appFile, nil
}

func (repo *AppFileRepositoryImpl) GetAppFilesByTask(taskID uint, fileType string) ([]model.AppFile, error) {
	var appFiles []model.AppFile
	if err := database.DB.Where("task_id = ? AND file_type = ?", taskID, fileType).Find(&appFiles).Error; err != nil {
		return nil, err
	}
	return appFiles, nil
}

func (repo *AppFileRepositoryImpl) GetAppFileByTask(taskID uint, fileType string) (*model.AppFile, error) {
	var appFile model.AppFile
	if err := database.DB.Where("task_id = ? AND file_type = ?", taskID, fileType).First(&appFile).Error; err != nil {
		return nil, err
	}
	return &appFile, nil
}
