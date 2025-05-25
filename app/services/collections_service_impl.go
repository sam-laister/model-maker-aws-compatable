package services

import (
	"github.com/Soup666/diss-api/model"
	repositories "github.com/Soup666/diss-api/repository"
)

type CollectionsServiceImpl struct {
	collectionsRepo repositories.CollectionsRepository
}

func NewCollectionsService(collectionsRepo repositories.CollectionsRepository) CollectionsService {
	return &CollectionsServiceImpl{collectionsRepo: collectionsRepo}
}

func (s *CollectionsServiceImpl) CreateCollection(collection *model.Collection) error {

	if err := s.collectionsRepo.CreateCollection(collection); err != nil {
		return err
	}
	return nil
}

func (s *CollectionsServiceImpl) GetCollection(collectionID uint) (*model.Collection, error) {
	report, err := s.collectionsRepo.GetCollectionByID(collectionID)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (s *CollectionsServiceImpl) GetCollections(userID uint) ([]model.Collection, error) {
	reports, err := s.collectionsRepo.GetCollectionsByUser(userID)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (s *CollectionsServiceImpl) ArchiveCollection(collectionID uint) error {
	err := s.collectionsRepo.ArchiveCollection(collectionID)
	if err != nil {
		return err
	}
	return nil
}

func (s *CollectionsServiceImpl) SaveCollection(collection *model.Collection) error {
	if err := s.collectionsRepo.SaveCollection(collection); err != nil {
		return err
	}
	return nil
}
