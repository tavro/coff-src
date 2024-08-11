package ast

import (
	"coff-src/src/coff/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&DefStatement{
				Token: token.Token{Type: token.DEF, Literal: "def"},
				Name: &Identifier{
					Token: token.Token{Type: token.ID, Literal: "thisVar"}
					Value: "thisVar"
				},
				Value: &Identifier{
					Token: token.Token{Type: token.ID, Literal: "thatVar"}
					Value: "thatVar"
				},
			},
		},
	}

	if program.String() != "def thisVar = thatVar;" {
		t.Errorf("program.String() is incorrect. got=%q", program.String())
	}
}