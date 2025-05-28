package mocks

import "github.com/stretchr/testify/mock"

// MockVisionService is a mock implementation of VisionService.
type MockVisionService struct {
	mock.Mock
}

func (m *MockVisionService) AnalyseImage(imagePath string, prompt string) (string, error) {
	args := m.Called(imagePath, prompt)
	if args.Get(0) != nil {
		return args.Get(0).(string), args.Error(1)
	}
	return "", args.Error(1)
}

func (m *MockVisionService) GenerateMessage(message string) (string, error) {
	args := m.Called(message)
	if args.Get(0) != nil {
		return args.Get(0).(string), args.Error(1)
	}
	return "", args.Error(1)
}
