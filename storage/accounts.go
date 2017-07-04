package storage

import (
	"errors"
	"github.com/bigblind/marvin/domain"
	"github.com/boltdb/bolt"
)

var accountFound = errors.New("Account found")

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

	if ab == nil {
		return act, domain.ErrAccountNotFound
	}

	return act, bytesToData(&act, ab)
}

func (s Store) GetAccountByEmail(email string) (domain.Account, error) {
	bucket, err := s.getOrCreateAccountsBucket()
	if err != nil {
		return domain.Account{}, err
	}

	act := domain.Account{}
	err = bucket.ForEach(func(_, ab []byte) error {
		err = bytesToData(&act, ab)
		if err == nil && act.Email == email {
			// Returning an error stops the iteration.
			// ForEach then returns the same error.
			// So we check for this error value later.
			return accountFound
		}
		return nil
	})

	if err == nil {
		return domain.Account{}, domain.ErrAccountNotFound
	} else if err != accountFound {
		return domain.Account{}, err
	}

	return act, nil
}

func (s Store) GetDefaultAccount() (act domain.Account, err error) {
	act, err = s.GetAccountByID("default")
	if err == nil {
		return
	} else if err == domain.ErrAccountNotFound {
		act, err = domain.NewAccount("", "")
		if err != nil {
			return
		}

		act.ID = "default"
		err = s.SaveAccount(act)
	}

	return
}

func (s Store) getOrCreateAccountsBucket() (*bolt.Bucket, error) {
	return s.tx.CreateBucketIfNotExists([]byte("accounts"))
}
