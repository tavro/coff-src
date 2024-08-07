package token

type TokenType string

type Token struct {
	Type 	TokenType
	Literal string
}

const (
	INVALID 	= "INVALID"
	EOF 		= "EOF"

	ID			= "ID"
	INT			= "INT"

	ASSIGN		= "="
	PLUS		= "+"

	COMMA 		= ","
	SEMICOLON	= ";"

	LPAR 		= "("
	RPAR 		= ")"
	LBRA		= "{"
	RBRA		= "}"

	FUN			= "FUN"
	LET			= "LET"
)