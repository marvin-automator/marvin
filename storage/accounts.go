package storage

import (
	"github.com/bigblind/marvin/domain"
	"github.com/boltdb/bolt"
)

func (s Store) SaveAccount(acct domain.Account) error {
	bucket, err := s.getOrCreateAccountsBucket()
	if err != nil {
		return err
	}

	ab, err := dataToBytes(acct)
	if err != nil {
		return err
	}

	return bucket.Put([]byte(acct.ID), ab)
}

func (s Store) GetAccountByID(aid string) (domain.Account, error) {
	bucket, err := s.getOrCreateAccountsBucket()
	if err != nil {
		return domain.Account{}, err
	}

	act := domain.Account{}
	ab := bucket.Get([]byte(aid))
	return act, bytesToData(*act, ab)
}

func (s Store) getOrCreateAccountsBucket() (*bolt.Bucket, error) {
	return s.tx.CreateBucketIfNotExists([]byte("accounts"))
}
