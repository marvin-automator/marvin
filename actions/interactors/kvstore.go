package interactors

import "github.com/marvin-automator/marvin/actions/domain"

// KVStoreStore persists values held by a KVStore
// Sorry about the slightly awkward name :).
// KVStores are stored in different buckets. KVStores within a bucket must have a unique ID.
type KVStoreStore interface {
	GetKVStore(bucket, ID string) (domain.KVStore, error)
	DeleteKVStore(bucket, ID string) error
}
