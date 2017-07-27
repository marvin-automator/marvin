package storage

import (
	"github.com/boltdb/bolt"
)

// A Store is the way you access stored data. When you know you're not going to need to write data, open a read-only
// store. There can only be one writable store open at a time, but multiple read-only stores.
type boltStore struct {
	tx *bolt.Tx
	writers int
}

// NewStore creates a new Store
func NewStore() Store {
	return boltStore{nil, 0}
}

func (bs *boltStore) beginWrite() error {
	var err error
	if bs.writers == 0 {
		bs.tx.Rollback()
		bs.tx, err = db.Begin(true)
	}
	bs.writers += 1
	return err
}

func (bs *boltStore) endWrite() error {
	var err error
	bs.writers -= 1

	if bs.writers == 1 {
		err = bs.tx.Commit()
		if err != nil {
			return err
		}
		bs.tx, err = db.Begin(false)
	}
	return err
}

func (bs *boltStore) newBucket(path []string) boltBucket {
	return boltBucket{bs, path}
}

func (bs *boltStore) Bucket(name string) (Bucket, error) {
	bs.beginWrite()
	_, err := bs.tx.CreateBucketIfNotExists([]byte(name))
	bs.endWrite()

	if err != nil {
		return boltBucket{}, err
	}
	return bs.newBucket([]string{name}), nil
}

func (bs *boltStore) DeleteBucket(name string) error {
	bs.beginWrite()
	err := bs.tx.DeleteBucket([]byte(name))
	bs.endWrite()
	return err
}

// DecodeBytes converts an array of bytes to a struct
// This should only be used by domain-specific Store implementations
func (s boltStore) DecodeBytes(d interface{}, b []byte) error {
	return bytesToData(d, b)
}

// EncodeBytes converts a struct that gets passd in to bytes
// This should only be used by domain-specific Store implementations
func (s boltStore) EncodeBytes(d interface{}) ([]byte, error) {
	return dataToBytes(d)
}


// GetBucketFromPath traverses the tree of buckets, to get the bucket at the given path.
func (bs *boltStore) GetBucketFromPath(path... string) (Bucket, error) {
	_, err := bs.getBoltBucketFromPath(path)
	return bs.newBucket(path), err
}

func (bs *boltStore) getBoltBucketFromPath(path []string) (*bolt.Bucket, error) {
	var current *bolt.Bucket
	for i, n := range path {
		if i == 0 {
			current = bs.tx.Bucket([]byte(n))
		} else {
			current = current.Bucket([]byte(n))
		}
		if current == nil {
			return current, NotFoundError
		}
	}
	return current, nil
}

func (bs *boltStore) Close() error {
	if bs.tx.Writable() {
		return bs.tx.Commit()
	}
	return bs.tx.Rollback()
}

type BoltBucket struct {
	store boltStore
	path []string
}

type boltBucket struct {
	bs *boltStore
	path []string
}

// bult returns the bolt Bucket corresponding to this bucket
func (bb *boltBucket) bolt() (*bolt.Bucket, error) {
	return bb.bs.getBoltBucketFromPath(bb.path)
}

// Bucket gets the child bucket with the given bucket, creating it
// if it doesn't exist.
func (bb *boltBucket) Bucket(name string) (Bucket, error) {
	b, err := bb.bolt()
	if err != nil {
		return boltBucket{}, err
	}

	_, err = b.CreateBucketIfNotExists([]byte(name))
	return bb.bs.newBucket(append(bb.path, name)), err
}

func (bb *boltBucket) DeleteBucket(name string) error {
	bb.bs.beginWrite()
	defer bb.bs.endWrite()

	b, err := bb.bs.getBoltBucketFromPath(bb.path)
	if err != nil {
		return err
	}

	return b.DeleteBucket([]byte(name))
}

func (bb *boltBucket) Get(key string, value interface{}) error {
	b, err := bb.bolt()
	if err != nil {
		return err
	}

	return bb.bs.DecodeBytes(value, b.Get([]byte(key)))
}

func (bb *boltBucket) Put(key string, value interface{}) error {
	b, err := bb.bolt()
	if err != nil {
		return err
	}

	bv, err := bb.bs.EncodeBytes(value)
	if err != nil {
		return err
	}

	return b.Put([]byte(key), bv)
}

func (bb *boltBucket) Each(f func(key string) error) error {
	b, err := bb.bolt()
	if err != nil {
		return err
	}

	return b.ForEach(func(key, value []byte) error {
		if value != nil {
			return f(string(key))
		}
		return nil
	})
}

func (bb *boltBucket) Delete(key string) error {
	b, err := bb.bolt()
	if err != nil {
		return err
	}

	return b.Delete([]byte(key))
}
