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
