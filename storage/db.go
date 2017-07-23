package storage

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
	"path"
	"errors"
)

var db *bolt.DB

// Setup initializes the database system, creating a database file
// called "marvin.db" in the current working directory.
func Setup() {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	SetupInFile(path.Join(wd, "marvin.db"))
}

// SetupTestDB initializes the database system, creating a database file in
// the os's temporary files directory.
func SetupTestDB() {
	tempDir := os.TempDir()
	SetupInFile(path.Join(tempDir, "marvin.db"))
}

// SetupInFile initializes the database system in a file and location specified by fname.
func SetupInFile(fname string) {
	var err error

	db, err = bolt.Open(fname, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// CloseDB closes the database.
func CloseDB() {
	err := db.Close()
	if err != nil {
		panic(err)
	}
}

// DeleteDBFile deletes the file associated with the currently configured database.
func DeleteDBFile() {
	path := db.Path()
	err := db.Close()
	if err != nil {
		panic(err)
	}
	err = os.Remove(path)
	if err != nil {
		panic(err)
	}
}

// A Store is the way you access stored data. When you know you're not going to need to write data, open a read-only
// store. There can only be one writable store open at a time, but multiple read-only stores.
type Store struct {
	Tx *bolt.Tx
}

// newStore creates a new Store
func newStore(wrt bool) (s Store, err error) {
	s = Store{}
	s.Tx, err = db.Begin(wrt)
	if err != nil {
		s = Store{}
	}

	return
}

// NewWritableStore creates a new Writable store
// Note that changes you make are only saved to disk when you close it,
// and that changes will only be visible to stores created after this one is closed.
func NewWritableStore() (Store, error) {
	return newStore(true)
}

// NewReadOnlyStore creates a new read-only store
func NewReadOnlyStore() (Store, error) {
	return newStore(false)
}

// Close closes the store. Always call either this or RollBack when you're done.
func (s Store) Close() {
	if s.Writable() {
		err := s.Tx.Commit()
		if err != nil {
			panic(err)
		}
	} else {
		err := s.Tx.Rollback()
		if err != nil {
			panic(err)
		}
	}
}

// RollBack rolls back any changes made to the store, and close it.
func (s Store) RollBack() {
	err := s.Tx.Rollback()
	if err != nil {
		panic(err)
	}
}

// DecodeBytes converts an array of bytes to a struct
// This should only be used by domain-specific Store implementations
func (s Store) DecodeBytes(d interface{}, b []byte) error {
	return bytesToData(d, b)
}

// EncodeBytes converts a struct that gets passd in to bytes
// This should only be used by domain-specific Store implementations
func (s Store) EncodeBytes(d interface{}) ([]byte, error) {
	return dataToBytes(d)
}

// Writable returns whether the store can be written to.
func (s Store) Writable() bool {
	return s.Tx.Writable()
}

// CreateBucketIfNotExists creates a bucket if it doesn't exist. This function opens a temporary writable store if the store
// on which the method is called is read-only.
func (s Store) CreateBucketIfNotExists(name string) (*bolt.Bucket, error) {
	var err error
	if s.Tx.Writable() {
		return s.Tx.CreateBucketIfNotExists([]byte(name))
	}
	b := s.Tx.Bucket([]byte(name))
	if b == nil {
		err = s.Tx.DB().Update(func(tx *bolt.Tx) error {
			_, err2 := tx.CreateBucket([]byte(name))
			return err2
		})
		b = s.Tx.Bucket([]byte(name))
	}
	return b, err
}

// CreateBucketHierarchy creates a hierarchy of buckets, starting at the root bucket. Once a child bucket
// is encountered down the path, that doesn't yet exist, it is created, and any further buckets in the path are created.
// This function opens a temporary writable store if the store on which the method is called is read-only.
// Returns the final bucket in the path.
func (s Store) CreateBucketHierarchy(path... string) (*bolt.Bucket, error) {
	var err error
	if s.Tx.Writable() {
		return s.createBucketHierarchy(path, s.Tx)
	}
	err = s.Tx.DB().Update(func(tx *bolt.Tx) error {
		_, err2 := s.createBucketHierarchy(path, tx)
		return err2
	})
	b, err := s.getBucketFromPath(path, s.Tx)
	return b, err
}

func (s Store) createBucketHierarchy(path []string, tx *bolt.Tx) (*bolt.Bucket, error) {
	var current *bolt.Bucket
	var err error
	for i, n := range path {
		if i == 0 {
			current, err = tx.CreateBucketIfNotExists([]byte(n))
		} else {
			current, err = current.CreateBucketIfNotExists([]byte(n))
		}
		if err != nil {
			return current, err
		}
	}
	return current, nil
}

func (s Store) getBucketFromPath(path []string, tx *bolt.Tx) (*bolt.Bucket, error) {
	var current *bolt.Bucket
	for i, n := range path {
		if i == 0 {
			current = tx.Bucket([]byte(n))
		} else {
			current = current.Bucket([]byte(n))
		}
		if current == nil {
			return current, errors.New("Bucket not found " + n)
		}
	}
	return current, nil
}