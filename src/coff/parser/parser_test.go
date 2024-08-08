package parser

import (
	"testing"
	"coff-src/src/coff/ast"
	"coff-src/src/coff/lexer"
)

func TestDefStatements(t *testing.T) {
	input := `
	def x 4;
	def = 0;
	def 123456;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedId string
	} {
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testDefStatement(t, stmt, tt.expectedId) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testDefStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "def" {
		t.Errorf("s.TokenLiteral is not 'def'. got=%q", s.TokenLiteral())
		return false
	}

	defStmt, ok := s.(*ast.DefStatement)
	if !ok {
		t.Errorf("s is not *ast.DefStatement. got=%T", s)
		return false
	}

	if defStmt.Name.Value != name {
		t.Errorf("defStmt.Name.Value is not '%s'. got=%s", name, defStmt.Name.Value)
		return false
	}

	if defStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name is not '%s'. got=%s", name, defStmt.Name)
		return false
	}
	
	return true
}