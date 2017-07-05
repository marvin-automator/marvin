package config

import (
	"github.com/bigblind/marvin/config/domain"
	"github.com/stretchr/testify/mock"
)

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
