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
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
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

