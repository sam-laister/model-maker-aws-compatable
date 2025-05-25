package mocks

import (
	models "github.com/Soup666/diss-api/model"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserFromFirebaseUID(apiKey string) (*models.User, error) {
	args := m.Called(apiKey)

	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)

	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)

	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *MockUserRepository) GetUsers() ([]*models.User, error) {
	args := m.Called()

	if args.Get(0) != nil {
		return args.Get(0).([]*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) DeleteUser(user *models.User) error {
	args := m.Called(user)

	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}
