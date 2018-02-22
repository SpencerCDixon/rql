package lexer

import (
	"testing"

	"github.com/spencercdixon/rql/testutil"
	"github.com/spencercdixon/rql/token"
)

func TestNew(t *testing.T) {
	upper := "SELECT FIRSTNAME FROM USERS;"
	expected := "select firstname from users;"

	l := New(upper)
	testutil.Assert(t, l.input == expected, "downcases the input used for lexing")
}

func TestSelect(t *testing.T) {
	query := `SELECT FirstName, LastName FROM users`
	l := New(query)

	tokens := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.SELECT, "select"},
		{token.IDENT, "firstname"},
		{token.COMMA, ","},
		{token.IDENT, "lastname"},
		{token.FROM, "from"},
		{token.IDENT, "users"},
		{token.EOF, ""},
	}

	for _, tt := range tokens {
		tok := l.NextToken()
		testutil.Equals(t, tt.expectedType, tok.Type)
		testutil.Equals(t, tt.expectedLiteral, tok.Literal)
	}
}
