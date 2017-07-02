package storage

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
	"path"
)

var db *bolt.DB

func Setup() {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	SetupInFile(path.Join(wd, "marvin.db"))
}

func SetupTestDB() {
	tempDir := os.TempDir()
	SetupInFile(path.Join(tempDir, "marvin.db"))
}

func SetupInFile(fname string) {
	var err error

	db, err = bolt.Open(fname, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	db.Close()
}

func DeleteDBFile() {
	db.Close()
	os.Remove(db.Path())
}

// A Store is the way you access stored data. When you know you're not going to need to write data, open a read-only
// store. There can only be one writable store open at a time, but multiple read-only stores.
type Store struct {
	tx *bolt.Tx
}

func newStore(wrt bool) (s Store, err error) {
	s = Store{}
	s.tx, err = db.Begin(wrt)
	if err != nil {
		s = Store{}
	}

	return
}

// Creates a new Writable store
// Note that changes you make are only saved to disk when you close it,
// and that changes will only be visible to stores created after this one is closed.
func NewWritableStore() (Store, error) {
	return newStore(true)
}

// Creates a new read-only store
func NewReadOnlyStore() (Store, error) {
	return newStore(false)
}

// Closes the store. Always call either this or RollBack when you're done.
func (s Store) Close() {
	s.tx.Commit()
}

// roll back any changes made to the store, and close it.
func (s Store) RollBack() {
	s.tx.Rollback()
}
