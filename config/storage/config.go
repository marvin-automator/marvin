package storage

import (
	"github.com/marvin-automator/marvin/config/domain"
	"github.com/marvin-automator/marvin/storage"
)

// ConfigStore is an implementation of the ConfigStore interface that uses the storage package as a backend.
type ConfigStore struct {
	storage.Store
}

// NewConfigStore creates a new ConfigStore
func NewConfigStore(s storage.Store) ConfigStore {
	return ConfigStore{s}
}

// GetConfig returns the current Config object
func (s ConfigStore) GetConfig() (domain.Config, error) {
	c := domain.DefaultConfig

	bucket, err := s.configBucket()
	if err != nil {
		return c, err
	}

	err = bucket.Get("config", &c)
	if err != nil && err != storage.NotFoundError {
		return c, err
	}
	return c, nil
}

// SaveConfig sets the given config object to be the current one.
func (s ConfigStore) SaveConfig(c domain.Config) error {
	bucket, err := s.configBucket()
	if err != nil {
		return err
	}

	return bucket.Put("config", c)
}

func (s ConfigStore) configBucket() (storage.Bucket, error) {
	return s.Bucket("config")
}
