package services

import models "github.com/Soup666/modelmaker/model"

type UserService interface {
	GetUserFromFirebaseUID(apiKey string) (*models.User, error)
	UpdateUser(user *models.User) error
}
