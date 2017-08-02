package storage

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
	"path"
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
