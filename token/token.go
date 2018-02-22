// Package token contains all of the lexical RQL tokens used by the lexer.
package token

import "strings"

// Type is a human readable form of our Tokens.  It is less efficient than
// using an iota with ints but makes building the toy DB easier to work with and
// debug.
type Type string

// Token represents a lexical token.
type Token struct {
	// Type is our human readable type for this token.
	Type Type
	// Literal is the actual value used to derive what type this token is.
	Literal string
}

const (
	ILLEGAL Type = "ILLEGAL"
	EOF          = "EOF"

	// Identifier
	IDENT      Type = "IDENT"
	STRING_TOK      = "STRING_TOK"
	INT_TOK         = "INT_TOK"

	// Operators
	ASSIGN Type = "="
	// LT          = "<"
	// GT          = ">"

	// Delimiters
	COMMA     Type = ","
	SEMICOLON      = ";"
	LPAREN         = "("
	RPAREN         = ")"

	// Keywords
	AND    Type = "AND"
	CREATE      = "CREATE"
	DELETE      = "DELETE"
	FROM        = "FROM"
	INDEX       = "INDEX"
	INSERT      = "INSERT"
	INTO        = "INTO"
	ON          = "ON"
	SELECT      = "SELECT"
	TABLE       = "TABLE"
	UPDATE      = "UPDATE"
	VALUES      = "VALUES"
	WHERE       = "WHERE"

	// Column Types
	VARCHAR = "VARCHAR"
	INT     = "INT"
)

var keywords = map[string]Type{
	"and":     AND,
	"create":  CREATE,
	"delete":  DELETE,
	"from":    FROM,
	"index":   INDEX,
	"insert":  INSERT,
	"int":     INT,
	"into":    INTO,
	"on":      ON,
	"select":  SELECT,
	"table":   TABLE,
	"update":  UPDATE,
	"values":  VALUES,
	"varchar": VARCHAR,
	"where":   WHERE,
}

func LookupIdent(ident string) Type {
	// allow case insensitivity for keywords
	ident = strings.ToLower(ident)
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
