package storage

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/marvin-automator/marvin/storage"
	"github.com/marvin-automator/marvin/actions/domain"
	"time"
)

func TestSaveAndGetTemporaryInstance(t *testing.T) {
	storage.WithTestDB(func(dbs storage.Store) {
		s := NewTemporaryInstanceStore(dbs)
		tai := domain.NewTemporaryActionInstance("provider", "action")

		err := s.Save(tai)
		require.NoError(t, err)

		tai2, err := s.Get(tai.ID)
		require.NoError(t, err)

		require.Equal(t, tai, tai2)
	})
}


func TestDeleteInstancesOlderThan(t *testing.T) {
	storage.WithTestDB(func(dbs storage.Store) {
		s := NewTemporaryInstanceStore(dbs)

		oldtai1 := domain.NewTemporaryActionInstance("old", "action")
		oldtai1.Created = time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)
		oldtai2 := domain.NewTemporaryActionInstance("old", "action2")
		oldtai2.Created = time.Date(2017, 2, 1, 0, 0, 0, 0, time.UTC)
		newtai := domain.NewTemporaryActionInstance("new", "action3")
		newtai.Created = time.Date(2017, 4, 1, 0, 0, 0, 0, time.UTC)

		err := s.Save(oldtai1)
		require.NoError(t, err)
		err = s.Save(oldtai2)
		require.NoError(t, err)
		err = s.Save(newtai)
		require.NoError(t, err)

		deleteAfter := time.Date(2017, 3, 1, 0, 0, 0, 0, time.UTC)
		n, err := s.DeleteInstancesOlderThan(deleteAfter)
		require.NoError(t, err)

		require.Equal(t, 2, n) // The number of instances deleted
		_, err = s.Get(oldtai1.ID)
		require.Error(t, err)
		_, err = s.Get(oldtai2.ID)
		require.Error(t, err)
		newtai2, err := s.Get(newtai.ID)
		require.NoError(t, err)
		require.Equal(t, newtai, newtai2)
	})
}