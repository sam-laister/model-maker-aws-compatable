package repositories

import (
	"github.com/Soup666/diss-api/model"
)

type AppFileRepository interface {
	SaveAppFile(appFile *model.AppFile) (*model.AppFile, error)
}
