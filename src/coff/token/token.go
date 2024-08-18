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
	STR			= "STR"

	ASSIGN		= "="
	PLUS		= "+" // Change to ADD?
	MINUS		= "-" // Change to SUB?
	MULT		= "*"
	DIV			= "/"
	HASH		= "#"
	FAC			= "!"

	QUERY		= "?" // TODO: Add to switch
	MOD 		= "%" // TODO: Add to switch

	LT 			= "<"
	GT 			= ">"

	EQ			= "=="
	NOT_EQ		= "!="

	COMMA 		= ","
	SEMICOLON	= ";"
	COLON 		= ":"

	LPAR 		= "("
	RPAR 		= ")"
	LBRA		= "{"
	RBRA		= "}"
	LBRACK		= "["
	RBRACK		= "]"

	FUN			= "FUN"
	DEF			= "DEF"
	TRUE 		= "TRUE"
	FALSE		= "FALSE"
	IF			= "IF"
	ELSE		= "ELSE"
	RET			= "RET"
	
	IS			= "IS"	// TODO: Add to keywords
	NOT 		= "NOT" // TODO: Add to keywords
	NUL			= "NUL" // TODO: Add to keywords
	NIL			= "NIL" // TODO: Add to keywords
)

var keywords = map[string]TokenType {
	"fun": FUN,
	"def": DEF,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"ret": RET,
}

func LookupId(id string) TokenType {
	if tok, ok := keywords[id]; ok {
		return tok
	}
	return ID
}