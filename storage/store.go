package storage

import "errors"

var (
	NotFoundError = errors.New("Not found.")
)

type Store interface {
	// Retrieves a storage bucket, creating it if it doesn't exist.
	Bucket(name string) (Bucket, error)
	// DeleteBucket deletes the bucket with the given name if it exists. Otherwise, it's a noop.
	DeleteBucket(name string) error

	// GetBucketFromPath traverses the tree of buckets. If a bucket along the path doesn't exist, NotFoundError is returned.
	GetBucketFromPath(path ...string) (Bucket, error)

	// CreateBucketHierarchy acts like GetBucketFromPath, but creates the buckets that don't exist.
	CreateBucketHierarchy(path ...string) (Bucket, error)

	// Close frees up the resources that the store uses.
	Close() error
}

type Bucket interface {
	// Retrieves a sub-bucket, creating it if it doesn't exist.
	Bucket(name string) (Bucket, error)
	// DeleteBucket deletes the bucket with the given name if it exists. Otherwise, it's a noop.
	DeleteBucket(name string) error

	// Retrieves the value associated with the given key.
	// The value passed in, should be a pointer to a value
	// of the type that is expected. That value will be replaced with the value from the store.
	// If there's no value for the given key; the pointer will be unaltered, and NotFoundError will be returned
	Get(key string, value interface{}) error

	// Stores the value at the given key in the bucket.
	Put(key string, value interface{}) error

	// Each calls f on every key in the bucket.
	// If f returns an error, iteration is stopped, and the error is returned.
	Each(f func(key string) error) error

	// Delete the value at the given key
	Delete(key string) error
}
