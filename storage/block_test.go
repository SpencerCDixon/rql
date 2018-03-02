package storage

import (
	"testing"

	"github.com/spencercdixon/rql/testutil"
)

func TestBlockEquals(t *testing.T) {
	b1 := NewBlock("users.tbl", 0)
	b2 := NewBlock("users.tbl", 1)
	b3 := NewBlock("users.tbl", 0)
	b4 := NewBlock("other.tbl", 0)

	testutil.Assert(t, b1.Equals(b3), "blocks are equal")
	testutil.Assert(t, !b1.Equals(b2), "blocks are not equal")
	testutil.Assert(t, !b1.Equals(b4), "blocks are not equal")
}

func TestBlockString(t *testing.T) {
	b1 := NewBlock("users.tbl", 0)
	testutil.Equals(t, b1.String(), "[file users.tbl, block 0]")
}
