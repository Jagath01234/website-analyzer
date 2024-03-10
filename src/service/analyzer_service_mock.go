package service

import (
	"github.com/stretchr/testify/mock"
	"website-analyzer/src/entity"
)

type MockAnalyzerService struct {
	mock.Mock
}

func (m *MockAnalyzerService) AnalyzeWebsiteContent(analysis entity.Analysis) (entity.Analysis, error) {
	args := m.Called(analysis)
	return args.Get(0).(entity.Analysis), args.Error(1)
}
