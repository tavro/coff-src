package lexer

import (
	"testing"
	"coff-src/src/coff/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType	token.TokenType
		expectedLiteral string
	} {
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAR, "("},
		{token.RPAR, ")"},
		{token.LBRA, "{"},
		{token.RBRA, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - incorrect token type. expected=%q, got=%q",
					 i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - incorrect literal. expected=%q, got=%q",
					 i, tt.expectedLiteral, tok.Literal)
		}
	}
}