package repositories

import (
	models "github.com/Soup666/modelmaker/model"
)

type ReportsRepository interface {
	GetReportsByUser(userID uint) ([]models.Report, error)
	GetReportByID(reportID uint) (*models.Report, error)
	CreateReport(report *models.Report) error
	SaveReport(report *models.Report) error
	ArchiveReport(reportID uint) error
}
