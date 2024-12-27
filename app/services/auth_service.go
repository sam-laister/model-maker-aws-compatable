package services

import (
	"errors"
	"log"

	"firebase.google.com/go/v4/auth"
	"gorm.io/gorm"

	database "github.com/Soup666/diss-api/database"
	model "github.com/Soup666/diss-api/model"
)

// AuthService provides authentication services
type AuthService struct {
	FireAuth *auth.Client
}

// Login authenticates a user with the provided credentials and returns a Firebase custom token
func (s *AuthService) Login(authToken string) (*model.User, error) {

	var user model.User
	err := database.DB.Where("firebase_uid = ?", authToken).First(&user).Error

	// Get the user from the database
	if err != nil {
		if err == gorm.ErrRecordNotFound {

			// First time - create user
			result := database.DB.Create(&model.User{FirebaseUid: authToken})
			if result.Error != nil {
				log.Printf("failed to insert user into database: %+v", result.Error)
				return nil, errors.New("internal server error")
			}
			err := database.DB.Where("firebase_uid = ?", authToken).First(&user).Error
			if err != nil {
				log.Printf("failed to get user by email from database: %v", err)
				return nil, errors.New("internal server error")
			}
			return &user, nil
		}
		log.Printf("failed to get user by email from database: %v", err)
		return nil, errors.New("internal server error")
	}

	return &user, nil
}
