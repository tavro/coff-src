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

	l.skipWhitespace()

	switch l.currChar {
	case '=':
		if l.peekChar() == '=' {
			currChar := l.currChar
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(currChar) + string(l.currChar)}
		} else {
			tok = newToken(token.ASSIGN, l.currChar)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.currChar)
	case ':':
		tok = newToken(token.COLON, l.currChar)
	case '(':
		tok = newToken(token.LPAR, l.currChar)
	case ')':
		tok = newToken(token.RPAR, l.currChar)
	case '{':
		tok = newToken(token.LBRA, l.currChar)
	case '}':
		tok = newToken(token.RBRA, l.currChar)
	case '[':
		tok = newToken(token.LBRACK, l.currChar)
	case ']':
		tok = newToken(token.RBRACK, l.currChar)
	case ',':
		tok = newToken(token.COMMA, l.currChar)
	case '+':
		tok = newToken(token.PLUS, l.currChar)
	case '-':
		tok = newToken(token.MINUS, l.currChar)
	case '*':
		tok = newToken(token.MULT, l.currChar)
	case '/':
		tok = newToken(token.DIV, l.currChar)
	case '#':
		tok = newToken(token.HASH, l.currChar)
	case '!':
		if l.peekChar() == '=' {
			currChar := l.currChar
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(currChar) + string(l.currChar)}
		} else {
			tok = newToken(token.FAC, l.currChar)
		}
	case '<':
		tok = newToken(token.LT, l.currChar)
	case '>':
		tok = newToken(token.GT, l.currChar)
	case '"':
		tok.Type = token.STR
		tok.Literal = l.readStr()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.currChar) {
			tok.Literal = l.readId()
			tok.Type = token.LookupId(tok.Literal)
			return tok
		} else if isDigit(l.currChar) {
			tok.Type = token.INT
			tok.Literal = l.readNum()
			return tok
		} else {
			tok = newToken(token.INVALID, l.currChar)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPos]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.currChar == ' ' || l.currChar == '\t' || l.currChar == '\n' || l.currChar == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readStr() string {
	position := l.pos + 1
	for {
		l.readChar()
		if l.currChar == '"' || l.currChar == 0 {
			break
		}
	}
	
	return l.input[position:l.pos]
}

func (l *Lexer) readNum() string {
	pos := l.pos
	for isDigit(l.currChar) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readId() string {
	pos := l.pos
	for isLetter(l.currChar) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}