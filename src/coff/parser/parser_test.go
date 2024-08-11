package parser

import (
	"fmt"
	"testing"
	"coff-src/src/coff/ast"
	"coff-src/src/coff/lexer"
)

func testIntLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntLiteral)
	if !ok {
		t.Errorf("il not *ast.IntLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value is not %d. got=%d", value, integ.Value)
		return false
	}
	
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input string
		expected string
	} {
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input string
		leftValue int64
		operator string
		rightValue int64
	} {
		{"4 + 4", 4, "+", 4},
		{"4 - 4", 4, "-", 4},
		{"4 * 4", 4, "*", 4},
		{"4 / 4", 4, "/", 4},
		{"4 > 4", 4, ">", 4},
		{"4 < 4", 4, "<", 4},
		{"4 == 4", 4, "==", 4},
		{"4 != 4", 4, "!=", 4},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}
		
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}
		if !testIntLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		
		if !testIntLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input string
		operator string
		intValue int64
	} {
		{"!4;", "!", 4},
		{"-10;", "-", 10},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		
		program := p.ParseProgram()
		checkParserErrors(t, p)
	
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}
		
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntLiteral(t, exp.Right, tt.intValue) {
			return
		}
	}
}

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