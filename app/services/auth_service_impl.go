package services

import (
	"context"
	"errors"

	"firebase.google.com/go/v4/auth"
	"github.com/Soup666/diss-api/model"
	repositories "github.com/Soup666/diss-api/repository"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	FireAuth *auth.Client
	DB       *gorm.DB
	userRepo repositories.UserRepository
}

func NewAuthService(FireAuth *auth.Client, DB *gorm.DB, userRepo repositories.UserRepository) AuthService {
	return &AuthServiceImpl{FireAuth: FireAuth, DB: DB, userRepo: userRepo}
}

func (s *AuthServiceImpl) ValidateToken(token string) (*auth.Token, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}

	authToken, err := s.FireAuth.VerifyIDToken(context.Background(), token)
	if err != nil {
		return nil, err
	}
	return authToken, nil
}

func (s *AuthServiceImpl) Verify(token string) (*model.User, error) {

	user, err := s.userRepo.GetUserFromFirebaseUID(token)

	if err == gorm.ErrRecordNotFound || user == nil {
		user := &model.User{
			FirebaseUid: token,
		}

		// Create blank user if not found
		err := s.userRepo.Create(user)

		if err != nil {
			return nil, errors.New("unable to verify user")
		}

		return user, nil
	} else if err != nil {
		return nil, errors.New("unable to verify user")
	}

	return user, nil
}

func (s *AuthServiceImpl) Unverify(user *model.User) error {
	err := s.userRepo.DeleteUser(user)
	if err != nil {
		return err
	}

	return nil
}
