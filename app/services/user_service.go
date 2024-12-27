package services

import (
	"errors"

	models "github.com/Soup666/diss-api/model"
	repositories "github.com/Soup666/diss-api/repositories"
)

// UserService is a struct that defines the service layer for users
type UserService struct {
}

// GetUserByAPIKey retrieves a user based on the API key
func (s *UserService) GetUserByAPIKey(apiKey string) (*models.User, error) {
	// Call the repository function to fetch user from the database
	user, err := repositories.GetUserByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
