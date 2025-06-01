package repositories

import (
	models "github.com/Soup666/modelmaker/model"
)

type UserAnalyticsRepository interface {
	GetAnalytics(userID uint) (*models.UserAnalytics, error)
}
