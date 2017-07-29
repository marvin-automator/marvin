package storage

import (
	"testing"
	"github.com/bigblind/marvin/storage"
	"github.com/bigblind/marvin/identityproviders/domain"
	"github.com/stretchr/testify/require"
)

func makeIdentity(id string) domain.Identity {
	return domain.Identity{
		ProviderID: id,
		Token: "Token",
		ImageURL: "http://example.com/image.jpg",
		Name: "John Doe",
		Subtext: "joe@example.com",
	}
}

func TestSaveAndGetIdentity(t *testing.T) {
	storage.WithTestDB(func(s storage.Store) {
		id := makeIdentity("007")
		is := NewIdentityStore(s)

		err := is.SaveIdentity("account", "testProvider", id)
		require.NoError(t, err)

		id2, err := is.GetIdentity("account", "testProvider", "007")
		require.NoError(t, err)

		require.Equal(t, id, id2)
	})
}

func TestIdentitiesForProvider(t *testing.T) {
	storage.WithTestDB(func(s storage.Store) {
		a1p1id1 := makeIdentity("007")
		a1p1id2 := makeIdentity("008")
		a1p2id1 := makeIdentity("009")
		a2p1id1 := makeIdentity("010")

		is := NewIdentityStore(s)

		err := is.SaveIdentity("account1", "p1", a1p1id1)
		require.NoError(t, err)
		err = is.SaveIdentity("account1", "p1", a1p1id2)
		require.NoError(t, err)
		err = is.SaveIdentity("account1", "p2", a1p2id1)
		require.NoError(t, err)
		err = is.SaveIdentity("account2", "p1", a2p1id1)
		require.NoError(t, err)

		ids, err := is.GetAccountIdentitiesForProvider("account1", "p1")
		require.NoError(t, err)

		require.Contains(t, ids, a1p1id1)
		require.Contains(t, ids, a1p1id2)
	})
}
