package services

import (
	"github.com/Soup666/diss-api/model"
)

type ReportsService interface {
	CreateReport(report *model.Report) error
	GetReport(reportID uint) (*model.Report, error)
	GetReports(userID uint) ([]model.Report, error)
	ArchiveReport(reportID uint) error
	SaveReport(report *model.Report) error
}
