package services

import (
	"github.com/Soup666/diss-api/model"
)

type CollectionsService interface {
	CreateCollection(collection *model.Collection) error
	GetCollection(collectionID uint) (*model.Collection, error)
	GetCollections(collections uint) ([]model.Collection, error)
	ArchiveCollection(collectionID uint) error
	SaveCollection(report *model.Collection) error
}
