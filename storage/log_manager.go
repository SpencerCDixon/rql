package storage

// LogManager is responsible for logging changes in our DBMS so they can be
// undone.  There is only ever one log file per DB.
type LogManager struct {
	// name of our log file
	filename   string
	currentBlk *Block
	currentPos int
}
