// Package lexer is used for turning SQL statement strings into tokens which can
// be used by the parser.
package lexer

import "strings"

/*
	Lexer supports five different token types:
		1. single character delimiteres, such as the comma
		2. integer constants, such as 123
		3. string constants, such as 'john'
		4. keywords, such as: select, from, and where
		5. identifiers (ident), such as: STUDENT, x, blahblerg
*/
type Lexer struct {
	input string
}

func New(input string) *Lexer {
	input = strings.ToLower(input)
	l := &Lexer{input: input}
	return l
}
