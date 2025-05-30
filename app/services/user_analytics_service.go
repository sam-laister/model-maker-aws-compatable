package services

import (
	"github.com/Soup666/modelmaker/model"
)

type UserAnalyticsService interface {
	GetAnalytics(userID uint) (*model.UserAnalytics, error)
}
