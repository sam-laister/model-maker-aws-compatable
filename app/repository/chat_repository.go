package repositories

import (
	"github.com/Soup666/modelmaker/model"
)

type ChatRepository interface {
	CreateChat(chat *model.ChatMessage) error
}
