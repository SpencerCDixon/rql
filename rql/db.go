package rql

// New creates a new database opening up a connection to the file where pages
// will be stored.
func New(path string) *Database {
	return &Database{path}
}

// Database is a rql database
type Database struct {
	Path string
}
