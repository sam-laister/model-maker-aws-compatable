package services

import (
	"github.com/Soup666/diss-api/model"
)

type AppFileService interface {
	Save(appFile *model.AppFile) (*model.AppFile, error)
}
