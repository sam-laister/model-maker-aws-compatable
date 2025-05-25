package mocks

import (
	"github.com/Soup666/diss-api/model"
	"github.com/stretchr/testify/mock"
)

type MockAppFileService struct {
	mock.Mock
}

func (m *MockAppFileService) Save(appFile *model.AppFile) (*model.AppFile, error) {
	args := m.Called(appFile)

	if args.Get(0) != nil {
		return args.Get(0).(*model.AppFile), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAppFileService) GetTaskFiles(taskID uint, fileType string) ([]model.AppFile, error) {
	args := m.Called(taskID, fileType)

	if args.Get(0) != nil {
		return args.Get(0).([]model.AppFile), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAppFileService) GetTaskFile(taskID uint, fileType string) (*model.AppFile, error) {
	args := m.Called(taskID, fileType)

	if args.Get(0) != nil {
		return args.Get(0).(*model.AppFile), args.Error(1)
	}
	return nil, args.Error(1)
}
