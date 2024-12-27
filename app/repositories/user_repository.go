package repositories

import (
	"github.com/Soup666/diss-api/database"
	models "github.com/Soup666/diss-api/model"
	"gorm.io/gorm"
)

// GetUserByAPIKey fetches a user by their API key from the database
func GetUserByAPIKey(apiKey string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("firebase_uid = ?", apiKey).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if not found, let the service layer handle this
		}
		return nil, err
	}
	return &user, nil
}
