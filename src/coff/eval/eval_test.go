package eval

import (
	"coff-src/src/coff/parser"
	"coff-src/src/coff/object"
	"coff-src/src/coff/lexer"
	"testing"
)

func TestEvalIntExpression(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	} {
		{"0", 0},
		{"4", 4},
		{"-1", -1},
		{"-4", -4},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	
	program := p.ParseProgram()
	
	return Eval(program)
}

func testIntObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Int)
	if !ok {
		t.Errorf("object is not an Int. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func TestEvalBoolExpression(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expected)
	}
}

func testBoolObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Bool)
	if !ok {
		t.Errorf("object is not a Bool. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

 	return true
}

func TestFacOperator(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expected)
	}
}

