package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Soup666/diss-api/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"firebase.google.com/go/v4/auth"
	models "github.com/Soup666/diss-api/model"
)

func TestAuthMiddleware(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	t.Run("Missing Authorization header", func(t *testing.T) {
		// Create a mock AuthService
		mockAuthService := new(services.MockAuthService)

		// Create the middleware
		middleware := AuthMiddleware(mockAuthService)

		// Create a test request
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		// Initialize Gin context
		c, _ := gin.CreateTestContext(recorder)
		c.Request = req

		// Call middleware
		middleware(c)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.JSONEq(t, `{"error":"Invalid Authorization header"}`, recorder.Body.String())
	})

	t.Run("Invalid token format", func(t *testing.T) {
		mockAuthService := new(services.MockAuthService)

		middleware := AuthMiddleware(mockAuthService)

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "InvalidToken")
		recorder := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(recorder)
		c.Request = req

		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.JSONEq(t, `{"error":"Invalid Authorization header"}`, recorder.Body.String())
	})

	t.Run("Token validation fails", func(t *testing.T) {
		mockAuthService := new(services.MockAuthService)
		mockAuthService.On("ValidateToken", "invalid-token").Return(nil, errors.New("Token validation failed"))

		middleware := AuthMiddleware(mockAuthService)

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		recorder := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(recorder)
		c.Request = req

		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.JSONEq(t, `{"error":"Token validation failed"}`, recorder.Body.String())
	})

	t.Run("User verification fails", func(t *testing.T) {
		mockAuthService := new(services.MockAuthService)
		mockAuthService.On("ValidateToken", "valid-token").Return(&auth.Token{UID: "123"}, nil)
		mockAuthService.On("Verify", "123").Return(nil, errors.New("unable to verify user"))

		middleware := AuthMiddleware(mockAuthService)

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		recorder := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(recorder)
		c.Request = req

		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		assert.JSONEq(t, `{"error":"Unable to verify user"}`, recorder.Body.String())
	})

	t.Run("Valid token and user", func(t *testing.T) {

		mockAuthService := new(services.MockAuthService)
		mockAuthService.On("ValidateToken", "valid-token").Return(&auth.Token{UID: "123"}, nil)
		mockAuthService.On("Verify", "123").Return(&models.User{Id: 1, FirebaseUid: "123", Email: ""}, nil)

		middleware := AuthMiddleware(mockAuthService)

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		recorder := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(recorder)
		c.Request = req

		middleware(c)

		user, _ := c.Get("user")
		token, _ := c.Get("token")
		c.JSON(http.StatusOK, gin.H{
			"user":  user,
			"token": token,
		})

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.JSONEq(t, `{
            "user": {"CreatedAt":"0001-01-01T00:00:00Z", "DeletedAt":null, "Email":"", "FirebaseUid":"123", "Id":1, "UpdatedAt":"0001-01-01T00:00:00Z"},
            "token": "123"
        }`, recorder.Body.String())
	})
}
