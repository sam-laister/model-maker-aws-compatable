package services

import (
	"github.com/Soup666/modelmaker/model"
)

type NotificationService interface {
	SendMessage(notification *model.Notification) (*model.Notification, error)
}
