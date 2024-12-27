package controller

import (
	"context"
	"net/http"
	"strings"

	services "github.com/Soup666/diss-api/services"

	"github.com/gin-gonic/gin"
)

// AuthController is the controller for handling authentication requests
type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService}
}

func (c *AuthController) Login(ctx *gin.Context) {

	// Extract API key from request header
	apiKey := ctx.GetHeader("Authorization")
	if apiKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "API key is missing"})
		return
	}

	// Remove "Bearer " if present
	apiKey = strings.TrimPrefix(apiKey, "Bearer ")

	authToken, err := c.authService.FireAuth.VerifyIDToken(context.Background(), apiKey)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid token"})
		return
	}

	// Register the new user and get a custom token for the user
	user, err := c.authService.Login(authToken.UID)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	// Return the custom token to the client
	ctx.JSON(200, gin.H{"user": user})
}
