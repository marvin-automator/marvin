package storage

import (
	"testing"
	"github.com/bigblind/marvin/storage"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strings"
)

func TestKVStoreStoreGetSetSameStore(t *testing.T) {
	storage.WithTestDB(t, func(dbs storage.Store) {
		s := NewKVStoreStore(dbs)

		kv, err := s.GetKVStore("bucket", "id")
		require.NoError(t, err)

		err = kv.Put("key", "value")
		require.NoError(t, err)
		retrieved, err := kv.Get("key")
		require.NoError(t, err)

		require.Equal(t, "value", retrieved)
	})
}

type testData struct {
	I int
	J string
	K struct {
		L []byte
		M *int
	}
}

func makeTestData() testData {
	bs := make([]byte, 10)
	rand.Read(bs)
	i := rand.Int()
	return testData{
		I: rand.Int(),
		J: strings.Repeat("foo", rand.Intn(5)),
		K:  struct {
		L []byte
		M *int
		} {
			L: bs,
			M: &i,
		},
	}
}

func TestKVStoreStoreGetSetSameStoreCustomValue(t *testing.T) {
	storage.WithTestDB(t, func(dbs storage.Store) {
		s := NewKVStoreStore(dbs)
		value := makeTestData()

		kv, err := s.GetKVStore("bucket", "id")
		require.NoError(t, err)

		err = kv.Put("key", value)
		require.NoError(t, err)
		retrieved, err := kv.Get("key")
		require.NoError(t, err)

		require.Equal(t, value, retrieved)
	})
}
