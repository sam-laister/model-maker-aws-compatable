package middleware

import (
	"net/http"
	"strings"

	"github.com/Soup666/diss-api/services"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks the header for a valid API token
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the `Authorization` header
		authHeader := c.GetHeader("Authorization")

		// Validate the header format (e.g., "Bearer <token>")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		authToken, err := authService.ValidateToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token validation failed"})
			c.Abort()
			return
		}

		user, err := authService.Verify(authToken.UID)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to verify user"})
			c.Abort()
			return
		}

		// Token is valid, proceed with the request
		c.Set("token", authToken.UID)
		c.Set("user", user)
		c.Next()
	}
}
