package storage

import (
	"github.com/bigblind/marvin/config/domain"
	"github.com/bigblind/marvin/storage"
	"github.com/boltdb/bolt"
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

	bucket, err := s.getOrCreateConfigBucket()
	if err != nil {
		return
	}

	b := bucket.Get([]byte("config"))
	if b != nil {
		err = s.DecodeBytes(&c, b)
	}

	return
}

// SaveConfig sets the given config object to be the current one.
func (s ConfigStore) SaveConfig(c domain.Config) error {
	bucket, err := s.getOrCreateConfigBucket()
	if err != nil {
		return err
	}

	b, err := s.EncodeBytes(c)
	if err != nil {
		return err
	}

	return bucket.Put([]byte("config"), b)
}

func (s ConfigStore) getOrCreateConfigBucket() (*bolt.Bucket, error) {
	return s.CreateBucketIfNotExists("config")
}
