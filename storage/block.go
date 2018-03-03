package storage

import "fmt"

// Block is unit of measurement for how pages manage disk access.
type Block struct {
	FileName string
	BlockNum int
}

// NewBlock returns a new block.
func NewBlock(filename string, blockNum int) *Block {
	return &Block{FileName: filename, BlockNum: blockNum}
}

// Equals let's us determine if two blocks are the same or not.  If they are
// from the same file and located at the same offset then they are the same
// block.
func (b *Block) Equals(other *Block) bool {
	return other.FileName == b.FileName && other.BlockNum == b.BlockNum
}

// String lets us pretty print blocks for debugging purposes: [file users.tbl, block 3]
func (b *Block) String() string {
	return fmt.Sprintf("[file %s, block %d]", b.FileName, b.BlockNum)
}
