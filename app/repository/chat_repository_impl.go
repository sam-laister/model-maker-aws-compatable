package repositories

import (
	"github.com/Soup666/diss-api/model"
	"gorm.io/gorm"
)

type ChatRepositoryImpl struct {
	DB *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &ChatRepositoryImpl{DB: db}
}

// CreateChat implements ChatRepository.
func (c *ChatRepositoryImpl) CreateChat(chat *model.ChatMessage) error {
	if err := c.DB.Create(chat).Error; err != nil {
		return err
	}
	return nil
}
