package repositories

import (
	models "github.com/Soup666/diss-api/model"
)

type UserRepository interface {
	GetUserFromFirebaseUID(apiKey string) (*models.User, error)
	Create(user *models.User) error
	UpdateUser(user *models.User) error
	GetUsers() ([]*models.User, error)
	DeleteUser(user *models.User) error
}
