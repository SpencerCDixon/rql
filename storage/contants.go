package storage

import "github.com/pkg/errors"

const (
	// BlockSize denotes the number of bytes in a Block.  Most DBMS use the
	// underlying OS's block size (generally 4KB) but for educational purposes
	// we're using an artifically low size to generate a lot of blocks.
	BlockSize = 400
	// IntSize represents how many bytes we will let the INT type be in our
	// database. We're using uint32 so this will be 4 bytes.
	IntSize = 4
)

var (
	// ErrPageFull is returned when trying to et an int or string to a page
	// that does not have enough room to set those primitives.
	ErrPageFull = errors.New("storage: not enough bytes available in this pages content")
)
