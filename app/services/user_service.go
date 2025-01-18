package services

import models "github.com/Soup666/diss-api/model"

type UserService interface {
	GetUserFromFirebaseUID(apiKey string) (*models.User, error)
	UpdateUser(user *models.User) error
}
