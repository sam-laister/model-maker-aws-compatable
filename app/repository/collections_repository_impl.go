package repositories

import (
	"github.com/Soup666/diss-api/database"
	models "github.com/Soup666/diss-api/model"
	"gorm.io/gorm"
)

type CollectionsRepositoryImpl struct {
	DB *gorm.DB
}

func NewCollectionsRepository(db *gorm.DB) CollectionsRepository {
	return &CollectionsRepositoryImpl{DB: db}
}

func (repo *CollectionsRepositoryImpl) GetCollectionsByUser(userID uint) ([]models.Collection, error) {
	var collections []models.Collection
	if err := database.DB.
		Preload("Tasks").
		Where("user_id = ?", userID).
		Find(&collections).Error; err != nil {
		return nil, err
	}
	return collections, nil
}

func (repo *CollectionsRepositoryImpl) GetCollectionByID(collectionID uint) (*models.Collection, error) {
	var collection models.Collection
	if err := database.DB.
		Preload("Tasks").
		Model(&models.Collection{}).
		Where("id = ?", collectionID).
		First(&collection).Error; err != nil {
		return nil, err
	}
	return &collection, nil
}

func (repo *CollectionsRepositoryImpl) CreateCollection(collection *models.Collection) error {
	if err := database.DB.
		Model(&models.Collection{}).
		Create(collection).Error; err != nil {
		return err
	}
	return nil
}

func (repo *CollectionsRepositoryImpl) SaveCollection(collection *models.Collection) error {
	if err := database.DB.Save(collection).Error; err != nil {
		return err
	}
	return nil
}

func (repo *CollectionsRepositoryImpl) ArchiveCollection(collectionID uint) error {
	if err := database.DB.Delete(&models.Collection{}, collectionID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *CollectionsRepositoryImpl) GetCollectionTasks(collectionID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := database.DB.Model(&models.Collection{}).
		Where("id = ?", collectionID).
		Association("Tasks").Find(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}
