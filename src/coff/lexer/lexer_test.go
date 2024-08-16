package lexer

import (
	"testing"
	"coff-src/src/coff/token"
)

func TestNextToken(t *testing.T) {
	input := `def four = 4;
	def zero = 0;
	
	def add = fun(x, y) {
		x + y;
	};
	
	def res = add(four, zero);
	"foobar"
	"foo bar"
	[1, 2];
	`

	tests := []struct {
		expectedType	token.TokenType
		expectedLiteral string
	} {
		{token.DEF, "def"},
		{token.ID, "four"},
		{token.ASSIGN, "="},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},
		{token.DEF, "def"},
		{token.ID, "zero"},
		{token.ASSIGN, "="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.DEF, "def"},
		{token.ID, "add"},
		{token.ASSIGN, "="},
		{token.FUN, "fun"},
		{token.LPAR, "("},
		{token.ID, "x"},
		{token.COMMA, ","},
		{token.ID, "y"},
		{token.RPAR, ")"},
		{token.LBRA, "{"},
		{token.ID, "x"},
		{token.PLUS, "+"},
		{token.ID, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRA, "}"},
		{token.SEMICOLON, ";"},
		{token.DEF, "def"},
		{token.ID, "res"},
		{token.ASSIGN, "="},
		{token.ID, "add"},
		{token.LPAR, "("},
		{token.ID, "four"},
		{token.COMMA, ","},
		{token.ID, "zero"},
		{token.RPAR, ")"},
		{token.SEMICOLON, ";"},
		{token.STR, "foobar"},
		{token.STR, "foo bar"},
		{token.LBRACK, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACK, "]"},
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