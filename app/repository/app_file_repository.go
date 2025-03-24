package repositories

import (
	"github.com/Soup666/diss-api/model"
)

type AppFileRepository interface {
	SaveAppFile(appFile *model.AppFile) (*model.AppFile, error)
	GetAppFilesByTask(taskID uint, fileType string) ([]model.AppFile, error)
	GetAppFileByTask(taskID uint, fileType string) (*model.AppFile, error)
}
