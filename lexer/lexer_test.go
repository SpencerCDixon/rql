package lexer

import (
	"testing"

	"github.com/spencercdixon/rql/testutil"
	"github.com/spencercdixon/rql/token"
)

func TestSelect(t *testing.T) {
	query := `SELECT FirstName, LastName FROM users`
	l := New(query)

	tokens := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.SELECT, "SELECT"},
		{token.IDENT, "FirstName"},
		{token.COMMA, ","},
		{token.IDENT, "LastName"},
		{token.FROM, "FROM"},
		{token.IDENT, "users"},
		{token.EOF, ""},
	}

	for _, tt := range tokens {
		tok := l.NextToken()
		testutil.Equals(t, tt.expectedType, tok.Type)
		testutil.Equals(t, tt.expectedLiteral, tok.Literal)
	}
}

func TestMilestoneOne(t *testing.T) {
	input := `
CREATE TABLE users (
  id int,
  name varchar(200),
  company varchar(100)
);
INSERT INTO users VALUES (1, 'Spencer Dixon', 'Rio');
SELECT name FROM users;
	`
	l := New(input)

	tokens := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.CREATE, "CREATE"},
		{token.TABLE, "TABLE"},
		{token.IDENT, "users"},
		{token.LPAREN, "("},
		{token.IDENT, "id"},
		{token.INT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "name"},
		{token.VARCHAR, "varchar"},
		{token.LPAREN, "("},
		{token.INT_TOK, "200"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.IDENT, "company"},
		{token.VARCHAR, "varchar"},
		{token.LPAREN, "("},
		{token.INT_TOK, "100"},
		{token.RPAREN, ")"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "users"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.INT_TOK, "1"},
		{token.COMMA, ","},
		{token.STRING_TOK, "Spencer Dixon"},
		{token.COMMA, ","},
		{token.STRING_TOK, "Rio"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.SELECT, "SELECT"},
		{token.IDENT, "name"},
		{token.FROM, "FROM"},
		{token.IDENT, "users"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	for _, tt := range tokens {
		tok := l.NextToken()
		t.Log(tt.expectedLiteral)
		testutil.Equals(t, tt.expectedType, tok.Type)
		testutil.Equals(t, tt.expectedLiteral, tok.Literal)
	}
}
