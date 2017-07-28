package storage

import (
	"github.com/bigblind/marvin/config/domain"
	"github.com/bigblind/marvin/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetConfigNoSaved(t *testing.T) {
	storage.WithTestDB(func(dbs storage.Store) {
		s := NewConfigStore(dbs)
		c, err := s.GetConfig()
		require.NoError(t, err)

		require.Equal(t, domain.DefaultConfig, c)
	})
}

func TestSaveAndGetConfig(t *testing.T) {
	storage.WithTestDB(func(dbs storage.Store) {
		s := NewConfigStore(dbs)
		c1 := domain.DefaultConfig
		c1.AccountsEnabled = true

		err := s.SaveConfig(c1)
		require.NoError(t, err)

		c2, err := s.GetConfig()
		require.NoError(t, err)

		require.Equal(t, c1, c2)
	})
}
