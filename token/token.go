// Package token defines constants representing the lexical tokens of the BF
// language
//
package token

import "strconv"

// Token is the set of lexical tokens of the Go programming language.
type Token int

// The list of tokens.
const (
	// Special tokens
	INVALID Token = iota
	EOF

	ADD // +
	SUB // -
	LSS // <
	GTR // >

	LBRACK // [

	RBRACK // ]
)

var tokens = [...]string{
	INVALID: "INVALID",
	EOF:     "EOF",

	ADD: "+",
	SUB: "-",

	LSS: "<",
	GTR: ">",

	LBRACK: "[",
	RBRACK: "]",
}

func (tok Token) String() string {
	s := ""

	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}

	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}

	return s
}
