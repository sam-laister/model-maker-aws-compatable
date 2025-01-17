package services

import (
	"firebase.google.com/go/v4/auth"
	"github.com/Soup666/diss-api/model"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	ValidateToken(token string) (*auth.Token, error)
	Verify(c *gin.Context) (*model.User, error)
}
