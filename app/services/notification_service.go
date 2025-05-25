package services

import (
	"github.com/Soup666/diss-api/model"
)

type NotificationService interface {
	SendMessage(notification *model.Notification) (*model.Notification, error)
}
