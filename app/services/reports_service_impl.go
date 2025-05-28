package services

import (
	"github.com/Soup666/diss-api/model"
	repositories "github.com/Soup666/diss-api/repository"
)

type ReportsServiceImpl struct {
	reportsRepo repositories.ReportsRepository
}

func NewReportsService(reportsRepo repositories.ReportsRepository) ReportsService {
	return &ReportsServiceImpl{reportsRepo: reportsRepo}
}

func (s *ReportsServiceImpl) CreateReport(report *model.Report) error {

	if err := s.reportsRepo.CreateReport(report); err != nil {
		return err
	}
	return nil
}

func (s *ReportsServiceImpl) GetReport(reportID uint) (*model.Report, error) {
	report, err := s.reportsRepo.GetReportByID(reportID)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (s *ReportsServiceImpl) GetReports(userID uint) ([]model.Report, error) {
	reports, err := s.reportsRepo.GetReportsByUser(userID)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (s *ReportsServiceImpl) ArchiveReport(reportID uint) error {
	err := s.reportsRepo.ArchiveReport(reportID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ReportsServiceImpl) SaveReport(report *model.Report) error {
	if err := s.reportsRepo.SaveReport(report); err != nil {
		return err
	}
	return nil
}
