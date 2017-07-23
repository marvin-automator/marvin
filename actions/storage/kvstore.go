package storage

import (
	"github.com/boltdb/bolt"
	"github.com/bigblind/marvin/storage"
	"github.com/bigblind/marvin/actions/domain"
	"github.com/bigblind/marvin/actions/interactors"
	"encoding/gob"
)

type kvStore struct {
	store storage.Store
	bucket *bolt.Bucket
}

// Get the value associated with a key
func (kv kvStore) Get(key string) (interface{}, error) {
	var v interface{}
	err := kv.store.DecodeBytes(&v, kv.bucket.Get([]byte(key)))
	return v, err
}

// Associate a value with the given key
func (kv kvStore) Put(key string, value interface{}) error {
	gob.Register(value)
	vb, err := kv.store.EncodeBytes(&value)
	if err != nil {
		return err
	}
	return kv.bucket.Put([]byte(key), vb)
}

// Delete the value at the given key.
func (kv kvStore) Delete(key string) error {
	return kv.bucket.Delete([]byte(key))
}

// Close closes the store, freeing up resouces.
func (kv kvStore) Close() {
	kv.store.Close()
}

// KVStore is an implementation of the domain.KVStoreStore interface, using the storage package
type KVStoreStore struct {
	store storage.Store
}

// NewKVStoreStore returns a new KVStoreStore.
func NewKVStoreStore(s storage.Store) interactors.KVStoreStore {
	return KVStoreStore{s}
}

// GetKVStore returns the KVStore in the given bucket with the given ID
func (k KVStoreStore) GetKVStore(bucket, ID string) (domain.KVStore, error){
	b, err := k.store.CreateBucketHierarchy("action_stores", bucket, ID)
	if err != nil {
		return nil, err
	}

	return kvStore{store: k.store, bucket: b}, nil
}

func (k KVStoreStore) DeleteKVStore(bucket, ID string) error {
	b, err := k.store.CreateBucketHierarchy("action_stores", bucket)
	if err != nil {
		return err
	}
	return b.DeleteBucket([]byte(ID))
}



