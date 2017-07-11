package config

import (
	"github.com/bigblind/marvin/config/domain"
	"github.com/stretchr/testify/mock"
)

// Mock ConfigStore implementation used in tests
type MockConfigStore struct {
	mock.Mock
}

func NewMockConfigStore() MockConfigStore {
	return MockConfigStore{mock.Mock{}}
}

func (m MockConfigStore) GetConfig() (domain.Config, error) {
	args := m.Called()
	return args.Get(0).(domain.Config), args.Error(1)
}
