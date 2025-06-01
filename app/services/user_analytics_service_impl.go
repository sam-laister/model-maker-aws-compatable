package services

import (
	"github.com/Soup666/modelmaker/model"
	repositories "github.com/Soup666/modelmaker/repository"
)

type UserAnalyticsServiceImpl struct {
	userAnalyticsRepo repositories.UserAnalyticsRepository
}

func NewUserAnalyticsService(userAnalyticsRepo repositories.UserAnalyticsRepository) UserAnalyticsService {
	return &UserAnalyticsServiceImpl{userAnalyticsRepo: userAnalyticsRepo}
}

func (s *UserAnalyticsServiceImpl) GetAnalytics(userID uint) (*model.UserAnalytics, error) {
	a, err := s.userAnalyticsRepo.GetAnalytics(userID)

	if err != nil {
		return nil, err
	}
	return a, nil
}
