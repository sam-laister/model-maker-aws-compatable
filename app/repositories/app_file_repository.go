package repositories

import (
	"github.com/Soup666/diss-api/database"
	"github.com/Soup666/diss-api/model"
)

// GetUserByAPIKey fetches a user by their API key from the database
func SaveAppFile(appFile *model.AppFile) (*model.AppFile, error) {
	if err := database.DB.Save(&appFile).Error; err != nil {
		return nil, err
	}
	return appFile, nil
}
