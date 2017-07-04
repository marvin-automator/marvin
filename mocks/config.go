package mocks

import "github.com/bigblind/marvin/domain"

type MockConfigStore struct {
	Config domain.Config
	Error error
}

func (m MockConfigStore) GetConfig() (domain.Config, error) {
	if m.Error != nil {
		return domain.Config{}, m.Error
	}
	return m.Config, nil
}
