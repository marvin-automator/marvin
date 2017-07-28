package storage

import (
	"github.com/bigblind/marvin/config/domain"
	"github.com/bigblind/marvin/storage"
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
func (s ConfigStore) GetConfig() (c domain.Config, err error) {
	c = domain.DefaultConfig

	bucket, err := s.configBucket()
	if err != nil {
		return
	}

	err = bucket.Get("config", &c)
	return
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
