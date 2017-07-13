package config

import (
	"github.com/bigblind/marvin/config/domain"
	"github.com/stretchr/testify/mock"
)

// MockConfigStore implementation used in tests
type MockConfigStore struct {
	*mock.Mock
}

// NewMockConfigStore returns a new MockConfigStore instance
func NewMockConfigStore() MockConfigStore {
	m := mock.Mock{}
	return MockConfigStore{&m}
}

// GetConfig returns the current Config instance.
func (m MockConfigStore) GetConfig() (domain.Config, error) {
	args := m.Called()
	return args.Get(0).(domain.Config), args.Error(1)
}
