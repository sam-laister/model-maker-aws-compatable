package services

import (
	"firebase.google.com/go/v4/auth"
	"github.com/Soup666/diss-api/model"
)

type AuthService interface {
	ValidateToken(token string) (*auth.Token, error)
	Verify(token string) (*model.User, error)
}
