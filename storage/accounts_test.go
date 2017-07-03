package storage

import (
	"testing"
	"github.com/bigblind/marvin/domain"
	"github.com/stretchr/testify/require"
	"bytes"
)

func TestSaveAndGetAccount(t *testing.T) {
	WithTestDB(t, func(s Store) {
		a1 := domain.Account{"042", "foo@example.com", []byte("nothashed")}
		err := s.SaveAccount(a1)
		require.NoError(t, err)

		a2, err := s.GetAccountByID("042")
		require.NoError(t, err)

		require.Equal(t, a1, a2)
	})
}

func TestGetAccountByEmailExists(t *testing.T) {
	WithTestDB(t, func(s Store) {
		a1 := domain.Account{"042", "foo@example.com", []byte("nothashed")}
		err := s.SaveAccount(a1)
		require.NoError(t, err)

		a2, err := s.GetAccountByEmail("foo@example.com")
		require.NoError(t, err)

		require.Equal(t, a1, a2, "The saved and retrieved accounts aren't equal.")

	})
}

func TestGetAccountByEmailDoesNotExist(t *testing.T) {
	WithTestDB(t, func(s Store) {
		a1 := domain.Account{"042", "foo@example.com", []byte("nothashed")}
		err := s.SaveAccount(a1)
		require.NoError(t, err)

		_, err = s.GetAccountByEmail("bar@example.com")
		require.EqualError(t, err, AccountNotFound.Error())
	})
}

func TestGetDefaultAccount(t *testing.T) {
	WithTestDB(t, func(s Store) {
		act, err := s.GetDefaultAccount()
		require.NoError(err)

		require.Equal(t, "default", act.ID, "Didn't return the default user.")
	})
}



