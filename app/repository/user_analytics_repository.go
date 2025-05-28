package repositories

import (
	models "github.com/Soup666/diss-api/model"
)

type UserAnalyticsRepository interface {
	GetAnalytics(userID uint) (*models.UserAnalytics, error)
}
