package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Soup666/diss-api/auth"
	"github.com/Soup666/diss-api/model"

	"github.com/gin-gonic/gin"
)

// AuthController is the controller for handling authentication requests
type AuthController struct {
	authService *auth.AuthService
}

func NewAuthController(authService *auth.AuthService) *AuthController {
	return &AuthController{authService}
}

func (c *AuthController) Login(ctx *gin.Context) {

	reqToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid token"})
		return
	}

	reqToken = strings.TrimSpace(splitToken[1])

	// Get the email from the request body
	var registrationData struct {
		Email string `json:"email"`
	}

	authToken, err := c.authService.FireAuth.VerifyIDToken(context.Background(), reqToken)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid token"})
		return
	}

	if err := ctx.ShouldBindJSON(&registrationData); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if registrationData.Email == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Email is required"})
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

// Register handles the POST /register route and creates a new user with the provided credentials
func (c *AuthController) Register(ctx *gin.Context) {

	reqToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid token"})
		return
	}

	reqToken = strings.TrimSpace(splitToken[1])

	authToken, err := c.authService.FireAuth.VerifyIDToken(context.Background(), reqToken)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid token"})
		return
	}

	// Get the email from the request body
	var registrationData struct {
		Email string `json:"email"`
	}

	if err := ctx.ShouldBindJSON(&registrationData); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if registrationData.Email == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "Email are required"})
		return
	}

	user := model.User{
		Email:       registrationData.Email,
		FirebaseUid: authToken.UID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	fmt.Println(user)

	// Register the new user and get a custom token for the user
	newUser, err := c.authService.Register(user)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	// Return the custom token to the client
	ctx.JSON(200, gin.H{"user": newUser})
}
