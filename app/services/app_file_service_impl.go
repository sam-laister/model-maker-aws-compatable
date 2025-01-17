package services

import (
	models "github.com/Soup666/diss-api/model"
	repositories "github.com/Soup666/diss-api/repository"
)

type AppFileServiceImpl struct {
	appFileRepo repositories.AppFileRepository
}

func NewAppFileServiceFile(appFileRepo repositories.AppFileRepository) *AppFileServiceImpl {
	return &AppFileServiceImpl{appFileRepo: appFileRepo}
}

func (s *AppFileServiceImpl) Save(appFile *models.AppFile) (*models.AppFile, error) {
	appFile, err := s.appFileRepo.SaveAppFile(appFile)
	if err != nil {
		return nil, err
	}
	return appFile, nil
}
