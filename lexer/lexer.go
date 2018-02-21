// Package lexer is used for turning SQL statement strings into tokens which can
// be used by the parser.
package lexer

import "strings"

type Lexer struct {
	input string
}

func New(input string) *Lexer {
	input = strings.ToLower(input)
	l := &Lexer{input: input}
	return l
}
