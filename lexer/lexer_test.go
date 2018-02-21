package lexer

import (
	"testing"

	"github.com/spencercdixon/rql/testutil"
)

func TestNew(t *testing.T) {
	upper := "SELECT FIRSTNAME FROM USERS;"
	expected := "select firstname from users;"

	l := New(upper)
	testutil.Assert(t, l.input == expected, "downcases the input used for lexing")
}
