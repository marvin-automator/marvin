package storage

import (
	"github.com/bigblind/marvin/domain"
	"github.com/boltdb/bolt"
)

type ConfigStore interface {
	GetConfig() (domain.Config, error)
}

func (s Store) GetConfig() (c domain.Config, err error) {
	c = domain.DefaultConfig

	bucket, err := s.getOrCreateConfigBucket()
	if err != nil {
		return
	}

	b := bucket.Get([]byte("config"))
	if b != nil {
		err = bytesToData(&c, b)
	}

	return
}

func (s Store) SaveConfig(c domain.Config) (e error) {
	bucket, err := s.getOrCreateConfigBucket()
	if err != nil {
		return
	}

	b, err := dataToBytes(c)
	if err != nil {
		return
	}

	err = bucket.Put([]byte("config"), b)
	return
}

func (s Store) getOrCreateConfigBucket() (*bolt.Bucket, error) {
	return s.tx.CreateBucketIfNotExists([]byte("config"))
}

