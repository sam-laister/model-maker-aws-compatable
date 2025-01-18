package controller

import (
	"net/http"

	"github.com/Soup666/diss-api/model"
	models "github.com/Soup666/diss-api/model"
	services "github.com/Soup666/diss-api/services"

	"github.com/gin-gonic/gin"
)

// AuthController is the controller for handling authentication requests
type AuthController struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthController(authService services.AuthService, userService services.UserService) *AuthController {
	return &AuthController{authService, userService}
}

func (c *AuthController) Verify(ctx *gin.Context) {

	switch ctx.Request.Method {
	case http.MethodPost:

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

	case http.MethodPatch:
		user := ctx.MustGet("user").(*model.User)

		var userUpdate models.User
		if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		// Fields allowed to be updated
		user.Email = userUpdate.Email

		err := c.userService.UpdateUser(user)

		if err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"user": user})
	default:
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}
}
