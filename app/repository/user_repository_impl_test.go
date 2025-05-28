package repositories_test

import (
	"testing"

	database "github.com/Soup666/diss-api/database"
	"github.com/Soup666/diss-api/model"
	repositories "github.com/Soup666/diss-api/repository"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository(t *testing.T) {

	err := database.SetupTestDB(t)

	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
		return
	}

	repo := repositories.NewUserRepository(database.DB)

	// Test Empty
	users, err := repo.GetUsers()
	assert.NoError(t, err)
	assert.Empty(t, users)
	assert.Equal(t, 0, len(users))

	// Test Create
	user := &model.User{
		FirebaseUid: "test_firebase_uid",
		Email:       "test@example.com",
	}
	err = repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.Model.ID)

	// Test GetUserFromFirebaseUID
	fetchedUser, err := repo.GetUserFromFirebaseUID("test_firebase_uid")
	assert.NoError(t, err)
	assert.NotNil(t, fetchedUser)
	assert.Equal(t, user.Email, fetchedUser.Email)

	// Test GetUserFromFirebaseUID with non-existent UID
	nonExistentUser, err := repo.GetUserFromFirebaseUID("nope")
	assert.NoError(t, err)
	assert.Nil(t, nonExistentUser)

	// Test UpdateUser
	user.Email = "test2@example.com"
	err = repo.UpdateUser(user)
	assert.NoError(t, err)

	updatedUser, _ := repo.GetUserFromFirebaseUID("test_firebase_uid")
	assert.Equal(t, "test2@example.com", updatedUser.Email)
}
