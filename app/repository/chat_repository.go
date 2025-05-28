package repositories

import (
	"github.com/Soup666/diss-api/model"
)

type ChatRepository interface {
	CreateChat(chat *model.ChatMessage) error
}
