package repositories

import (
	models "github.com/Soup666/diss-api/model"
)

type CollectionsRepository interface {
	GetCollectionsByUser(userID uint) ([]models.Collection, error)
	GetCollectionByID(collectionID uint) (*models.Collection, error)
	CreateCollection(collection *models.Collection) error
	SaveCollection(collection *models.Collection) error
	ArchiveCollection(collectionID uint) error
	GetCollectionTasks(collectionID uint) ([]models.Task, error)
}
