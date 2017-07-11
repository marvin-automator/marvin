package storage

import "testing"

// WithTestDB sets up the storage package with a test database file, and calls the passed-in function.
// When that function returns, the test database is removed.
func WithTestDB(t *testing.T, f func(Store)) {
	SetupTestDB()
	s, err := NewWritableStore()
	if err != nil {
		t.Error(err)
	} else {
		f(s)
		s.Close()
	}
	DeleteDBFile()
}
