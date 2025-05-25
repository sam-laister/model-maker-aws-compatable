package mocks

import (
	models "github.com/Soup666/diss-api/model"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository.
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetArchivedTasks(userID uint) ([]*models.Task, error) {
	args := m.Called(userID)

	if args.Get(0) != nil {
		return args.Get(0).([]*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) UnarchiveTask(taskID uint) (*models.Task, error) {
	args := m.Called(taskID)

	if args.Get(0) != nil {
		return args.Get(0).(*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) GetUnarchivedTasks(userID uint) ([]*models.Task, error) {
	args := m.Called(userID)

	if args.Get(0) != nil {
		return args.Get(0).([]*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) GetTaskByID(taskID uint) (*models.Task, error) {
	args := m.Called(taskID)

	if args.Get(0) != nil {
		return args.Get(0).(*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) CreateTask(task *models.Task) error {
	args := m.Called(task)

	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *MockTaskRepository) SaveTask(task *models.Task) error {
	args := m.Called(task)

	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *MockTaskRepository) ArchiveTask(task uint) (*models.Task, error) {
	args := m.Called(task)

	if args.Get(0) != nil {
		return nil, args.Error(0)
	}
	return nil, args.Error(0)
}

func (m *MockTaskRepository) AddLog(taskID uint, log string) error {
	args := m.Called(taskID, log)

	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}
