package storage

import (
	"github.com/bigblind/marvin/config/domain"
	"github.com/bigblind/marvin/storage"
	"github.com/boltdb/bolt"
)

type ConfigStore struct {
	storage.Store
}

func NewConfigStore(s storage.Store) ConfigStore {
	return ConfigStore{s}
}

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

func (s ConfigStore) SaveConfig(c domain.Config) (e error) {
	bucket, err := s.getOrCreateConfigBucket()
	if err != nil {
		return
	}

	b, err := s.EncodeBytes(c)
	if err != nil {
		return
	}

	err = bucket.Put([]byte("config"), b)
	return
}

func (s ConfigStore) getOrCreateConfigBucket() (*bolt.Bucket, error) {
	return s.Tx.CreateBucketIfNotExists([]byte("config"))
}
