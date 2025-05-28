package services

import (
	"github.com/Soup666/diss-api/model"
)

type UserAnalyticsService interface {
	GetAnalytics(userID uint) (*model.UserAnalytics, error)
}
