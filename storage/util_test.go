package storage

import (
	"testing"

	"github.com/spencercdixon/rql/testutil"
)

func TestByteSize(t *testing.T) {
	testutil.Assert(t, ByteSizeForVal(2000) == 4, "four bytes for ints")
	testutil.Assert(
		t,
		ByteSizeForVal("hello") == 9,
		"extra 4 bytes for strings plus 1 byte per rune",
	)
}
