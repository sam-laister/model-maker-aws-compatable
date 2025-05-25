package mocks

import (
	"github.com/Soup666/diss-api/model"
	"github.com/stretchr/testify/mock"
)

type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) CreateChat(chat *model.ChatMessage) error {
	args := m.Called(chat)

	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}
