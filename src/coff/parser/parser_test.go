package parser

import (
	"testing"
	"coff-src/src/coff/ast"
	"coff-src/src/coff/lexer"
)

func TestIntLiteralExpression(t *testing.T) {
	input := "4;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
		program.Statements[0])
	}
	
	literal, ok := stmt.Expression.(*ast.IntLiteral)
	if !ok {
		t.Fatalf("exp is not *ast.IntLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 4 {
		t.Errorf("literal.Value is not %d. got=%d", 4, literal.Value)
	}
	if literal.TokenLiteral() != "4" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "4", literal.TokenLiteral())
	}
}

func TestDefStatements(t *testing.T) {
	input := `
	def x = 4;
	def y = 0;
	def foobar = 123456;
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

func TestRetStatements(t *testing.T) {
	input := `
	ret 4;
	ret 0;
	ret 123456;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		retStmt, ok := stmt.(*ast.RetStatement)
		if !ok {
			t.Errorf("stmt is not *ast.retStatement. got=%T", stmt)
			continue
		}
		if retStmt.TokenLiteral() != "ret" {
			t.Errorf("retStmt.TokenLiteral is not 'ret', got %q",
			retStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp is not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value is not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral is not %s. got=%s", "foobar", ident.TokenLiteral())
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