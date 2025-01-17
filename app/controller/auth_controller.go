package controller

import (
	services "github.com/Soup666/diss-api/services"

	"github.com/gin-gonic/gin"
)

// AuthController is the controller for handling authentication requests
type AuthController struct {
	authService *services.AuthServiceImpl
}

func NewAuthController(authService *services.AuthServiceImpl) *AuthController {
	return &AuthController{authService}
}

func (c *AuthController) Verify(ctx *gin.Context) {

	token, exists := ctx.Get("token")
	if !exists {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Token not found"})
		return
	}

	// Register the new user and get a custom token for the user
	user, err := c.authService.Verify(token.(string))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	// Return the custom token to the client
	ctx.JSON(200, gin.H{"user": user})
}
