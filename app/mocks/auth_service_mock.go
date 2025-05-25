package mocks

import (
	"firebase.google.com/go/v4/auth"
	models "github.com/Soup666/diss-api/model"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthService.
type MockAuthService struct {
	mock.Mock
}

// ValidateToken mocks the token validation method.
func (m *MockAuthService) ValidateToken(token string) (*auth.Token, error) {
	args := m.Called(token)

	if args.Get(0) != nil {
		return args.Get(0).(*auth.Token), args.Error(1)
	}
	return nil, args.Error(1)
}

// Verify mocks the user verification method.
func (m *MockAuthService) Verify(token string) (*models.User, error) {
	args := m.Called(token)

	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// Unverify mocks the user unverification method.
func (m *MockAuthService) Unverify(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
