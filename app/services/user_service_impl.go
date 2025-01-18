package services

import (
	"errors"

	models "github.com/Soup666/diss-api/model"
	repositories "github.com/Soup666/diss-api/repository"
)

type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{userRepo: userRepo}
}

func (s *UserServiceImpl) GetUserFromFirebaseUID(apiKey string) (*models.User, error) {

	if apiKey == "" {
		return nil, errors.New("api key is required")
	}

	user, err := s.userRepo.GetUserFromFirebaseUID(apiKey)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserServiceImpl) UpdateUser(user *models.User) error {
	return s.userRepo.UpdateUser(user)
}
