package lexer

import "coff-src/src/coff/token"

type Lexer struct {
	input		string
	pos			int
	readPos 	int
	currChar	byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.currChar = 0 // ASCII "NUL"
	} else {
		l.currChar = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.currChar {
	case '=':
		tok = newToken(token.ASSIGN, l.currChar)
	case ';':
		tok = newToken(token.SEMICOLON, l.currChar)
	case '(':
		tok = newToken(token.LPAR, l.currChar)
	case ')':
		tok = newToken(token.RPAR, l.currChar)
	case '{':
		tok = newToken(token.LBRA, l.currChar)
	case '}':
		tok = newToken(token.RBRA, l.currChar)
	case ',':
		tok = newToken(token.COMMA, l.currChar)
	case '+':
		tok = newToken(token.PLUS, l.currChar)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}