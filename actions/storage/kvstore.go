package storage

import (
	"github.com/bigblind/marvin/storage"
	"github.com/bigblind/marvin/actions/domain"
	"github.com/bigblind/marvin/actions/interactors"
	"encoding/gob"
	"bytes"
)

type kvStore struct {
	store storage.Store
	bucket storage.Bucket
}

// Get the value associated with a key
func (kv kvStore) Get(key string) (interface{}, error) {
	var v interface{}
	bytes := []byte{}
	err := kv.bucket.Get(key, &bytes)
	if err != nil {
		return nil, err
	}

	err = gobdecode(&v, bytes)
	return v, err
}

// Associate a value with the given key
func (kv kvStore) Put(key string, value interface{}) error {
	vb, err := gobEncode(value)
	if err != nil {
		return err
	}
	return kv.bucket.Put(key, vb)
}

// Delete the value at the given key.
func (kv kvStore) Delete(key string) error {
	return kv.bucket.Delete(key)
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
	b, err := k.getBucket(bucket, ID)
	if err != nil {
		return nil, err
	}

	return kvStore{store: k.store, bucket: b}, nil
}

func (k KVStoreStore) DeleteKVStore(bucket, ID string) error {
	return k.store.DeleteBucket("kv_store_" + bucket + "_" + ID)
}

func (k KVStoreStore) getBucket(bucket, id string) (storage.Bucket, error) {
	return k.store.Bucket("kv_stores_" + bucket + "_" + id)
}

func gobEncode(value interface{}) ([]byte, error) {
	gob.Register(value)
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	err := enc.Encode(&value)
	return buf.Bytes(), err
}

func gobdecode(valptr interface{}, encoded []byte) error {
	dec := gob.NewDecoder(bytes.NewBuffer(encoded))
	return dec.Decode(valptr)
}


