package storage

import (
	"github.com/marvin-automator/marvin/accounts/domain"
	"github.com/marvin-automator/marvin/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSaveAndGetAccount(t *testing.T) {
	storage.WithTestDB(func(dbs storage.Store) {
		s := NewAccountStore(dbs)
		a1 := domain.Account{"042", "foo@example.com", []byte("nothashed")}
		err := s.SaveAccount(a1)
		require.NoError(t, err)

		a2, err := s.GetAccountByID("042")
		require.NoError(t, err)

		require.Equal(t, a1, a2)
	})
}

func TestGetAccountByEmailExists(t *testing.T) {
	storage.WithTestDB(func(dbs storage.Store) {
		s := NewAccountStore(dbs)
		a1 := domain.Account{"042", "foo@example.com", []byte("nothashed")}
		err := s.SaveAccount(a1)
		require.NoError(t, err)

		a2, err := s.GetAccountByEmail("foo@example.com")
		require.NoError(t, err)

		require.Equal(t, a1, a2, "The saved and retrieved accounts aren't equal.")

	})
}

func TestGetAccountByEmailDoesNotExist(t *testing.T) {
	storage.WithTestDB(func(dbs storage.Store) {
		s := NewAccountStore(dbs)
		a1 := domain.Account{"042", "foo@example.com", []byte("nothashed")}
		err := s.SaveAccount(a1)
		require.NoError(t, err)

		_, err = s.GetAccountByEmail("bar@example.com")
		require.EqualError(t, err, domain.ErrAccountNotFound.Error())
	})
}

func TestGetDefaultAccount(t *testing.T) {
	storage.WithTestDB(func(dbs storage.Store) {
		s := NewAccountStore(dbs)
		act, err := s.GetDefaultAccount()
		require.NoError(t, err)

		require.Equal(t, "default", act.ID, "Didn't return the default user.")
	})
}

func TestDeleteAccount(t *testing.T) {
	storage.WithTestDB(func(dbs storage.Store) {
		s := NewAccountStore(dbs)
		// Insert an account
		act := domain.Account{"042", "foo@example.com", []byte("nothashed")}
		err := s.SaveAccount(act)
		require.NoError(t, err)

		// Delete it
		err = s.DeleteAccount("042")
		require.NoError(t, err)

		// We shouldn't get an account back.
		_, err = s.GetAccountByID("042")
		require.EqualError(t, err, domain.ErrAccountNotFound.Error())
	})
}

func TestEachAccount(t *testing.T) {
	storage.WithTestDB(func(dbs storage.Store) {
		s := NewAccountStore(dbs)
		//insert some accounts
		acs := map[string]domain.Account{}
		acs["042"] = domain.Account{"042", "foo@example.com", []byte("nothashed")}
		acs["043"] = domain.Account{"043", "bar@example.com", []byte("nothashed")}
		acs["044"] = domain.Account{"044", "baz@example.com", []byte("nothashed")}
		err := s.SaveAccount(acs["042"])
		require.NoError(t, err)
		err = s.SaveAccount(acs["043"])
		require.NoError(t, err)
		err = s.SaveAccount(acs["044"])
		require.NoError(t, err)

		err = s.EachAccount(func(a domain.Account) error {
			delete(acs, a.ID)
			return nil
		})

		// It shouldn't have returned an error
		require.NoError(t, err)
		// It should have been called with all the accounts, removing all of them from the map
		require.Equal(t, 0, len(acs))
	})
}
