package services_test

import (
	"errors"
	"testing"

	"github.com/Soup666/diss-api/mocks"
	models "github.com/Soup666/diss-api/model"
	"github.com/Soup666/diss-api/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserService(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	userService := services.NewUserService(mockUserRepository)

	t.Run("GetUserFromFirebaseUID", func(t *testing.T) {
		apiKey := "valid_api_key"
		expectedUser := &models.User{
			Model:       gorm.Model{ID: 1},
			Email:       "example@example.com",
			FirebaseUid: "valid_api_key",
		}

		mockUserRepository.On("GetUserFromFirebaseUID", apiKey).Return(expectedUser, nil)
		mockUserRepository.On("GetUserFromFirebaseUID", "invalid_api_key").Return(nil, nil)
		mockUserRepository.On("GetUserFromFirebaseUID", "").Return(nil, nil)
		mockUserRepository.On("GetUserFromFirebaseUID", "error").Return(nil, errors.New("error"))

		fetchedUser, err := userService.GetUserFromFirebaseUID(apiKey)

		mockUserRepository.AssertCalled(t, "GetUserFromFirebaseUID", apiKey)

		assert.NoError(t, err)
		assert.NotNil(t, fetchedUser)
		assert.Equal(t, expectedUser.ID, fetchedUser.ID)
		assert.Equal(t, expectedUser.Email, fetchedUser.Email)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		user := &models.User{
			Model:       gorm.Model{ID: 1},
			Email:       "example@example.com",
			FirebaseUid: "valid_api_key",
		}

		mockUserRepository.On("UpdateUser", user).Return(nil)
		mockUserRepository.On("UpdateUser", nil).Return(errors.New("error"))
		mockUserRepository.On("UpdateUser", &models.User{}).Return(errors.New("error"))
		mockUserRepository.On("UpdateUser", &models.User{Model: gorm.Model{ID: 1}}).Return(nil)

		err := userService.UpdateUser(user)

		mockUserRepository.AssertCalled(t, "UpdateUser", user)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotNil(t, user.ID)
		assert.Equal(t, user.Email, "example@example.com")
	})
}
