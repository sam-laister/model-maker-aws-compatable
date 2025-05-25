package mocks

import (
	"github.com/Soup666/diss-api/model"
	"github.com/stretchr/testify/mock"
)

type MockNotificationService struct {
	mock.Mock
}

// ValidateToken mocks the token validation method.
func (m *MockNotificationService) SendMessage(notification *model.Notification) (*model.Notification, error) {
	args := m.Called(notification)

	if args.Get(0) != nil {
		return args.Get(0).(*model.Notification), args.Error(1)
	}
	return nil, args.Error(1)
}
