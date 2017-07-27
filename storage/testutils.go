package storage

// WithTestDB sets up the storage package with a test database file, and calls the passed-in function.
// When that function returns, the test database is removed.
func WithTestDB(f func(Store)) {
	SetupTestDB()
	s := NewStore()
	f(s)
	s.Close()
	DeleteDBFile()
}
