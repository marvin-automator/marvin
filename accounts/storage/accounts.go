package storage

import (
	"errors"
	"github.com/bigblind/marvin/accounts/domain"
	"github.com/bigblind/marvin/storage"
)

var accountFound = errors.New("account found") // nolint

// AccountStore is an implementation of the AccountStore interface that uses the storage package
type AccountStore struct {
	storage.Store
}

// NewAccountStore returns a new AccountStore
func NewAccountStore(s storage.Store) AccountStore {
	return AccountStore{s}
}

// SaveAccount saves an account to the store
func (s AccountStore) SaveAccount(acct domain.Account) error {
	bucket, err := s.accountsBucket()
	if err != nil {
		return err
	}

	return bucket.Put(acct.ID, acct)
}

// GetAccountByID returns the account with the given ID. If no such account exists, domain.ErrAccountNotFound is returned as the error.
func (s AccountStore) GetAccountByID(aid string) (domain.Account, error) {
	bucket, err := s.accountsBucket()
	if err != nil {
		return domain.Account{}, err
	}

	act := domain.Account{}
	err = bucket.Get(aid, &act)

	return act, err
}

// GetAccountByEmail returns the account with the given email adress. If no such account exists, domain.ErrAccountNotFound is returned as the error.
func (s AccountStore) GetAccountByEmail(email string) (domain.Account, error) {
	bucket, err := s.accountsBucket()
	if err != nil {
		return domain.Account{}, err
	}

	act := domain.Account{}
	err = bucket.Each(func(aid string) error {
		err = bucket.Get(aid, &act)
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

// GetDefaultAccount returns the default account instance.
func (s AccountStore) GetDefaultAccount() (act domain.Account, err error) {
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

// DeleteAccount deletes the account with the given ID.
func (s AccountStore) DeleteAccount(aid string) error {
	b, err := s.accountsBucket()
	if err != nil {
		return err
	}

	return b.Delete(aid)
}

// EachAccount executes f with each account in the system. If the database encounters an error,
// it is returned. If f returns an error, iteration is stopped and the error is returned.
func (s AccountStore) EachAccount(f func(a domain.Account) error) error {
	b, err := s.accountsBucket()
	if err != nil {
		return err
	}

	return b.Each(func(aid string) error {
		act := domain.Account{}
		err := b.Get(aid, &act)
		if err != nil {
			return err
		}

		return f(act)
	})
}

func (s AccountStore) accountsBucket() (storage.Bucket, error) {
	return s.Bucket("accounts")
}
