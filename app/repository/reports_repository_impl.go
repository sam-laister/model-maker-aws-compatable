package repositories

import (
	"github.com/Soup666/diss-api/database"
	models "github.com/Soup666/diss-api/model"
	"gorm.io/gorm"
)

type ReportsRepositoryImpl struct {
	DB *gorm.DB
}

func NewReportsRepository(db *gorm.DB) ReportsRepository {
	return &ReportsRepositoryImpl{DB: db}
}

func (repo *ReportsRepositoryImpl) GetReportsByUser(userID uint) ([]models.Report, error) {
	var reports []models.Report
	if err := database.DB.Model(&models.Report{}).Where("user_id = ?", userID).Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (repo *ReportsRepositoryImpl) GetReportByID(reportID uint) (*models.Report, error) {
	var report models.Report
	if err := database.DB.Model(&models.Report{}).Where("id = ?", reportID).First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (repo *ReportsRepositoryImpl) CreateReport(report *models.Report) error {
	if err := database.DB.Model(&models.Report{}).Create(report).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ReportsRepositoryImpl) SaveReport(report *models.Report) error {
	if err := database.DB.Save(report).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ReportsRepositoryImpl) ArchiveReport(reportID uint) error {
	if err := database.DB.Delete(&models.Report{}, reportID).Error; err != nil {
		return err
	}

	return nil
}
