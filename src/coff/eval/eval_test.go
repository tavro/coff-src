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

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		expectedMessage string
	} {
		{
			"5 + true;",
			"type mismatch: INT + BOOL",
		},
		{
			"5 + true; 5;",
			"type mismatch: INT + BOOL",
		},
		{
			"-true",
			"unknown operator: -BOOL",
		},
		{
			"true + false;",
			"unknown operator: BOOL + BOOL",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOL + BOOL",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOL + BOOL",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					ret true + false;
				}
				ret 1;
			}
			`,
			"unknown operator: BOOL + BOOL",
		},
	}
		
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestRetStatements(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	} {
		{"ret 10;", 10},
		{"ret 10; 9;", 10},
		{"ret 2 * 5; 9;", 10},
		{"9; ret 2 * 5; 9;", 10},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    ret 10;
  }

  ret 1;
}
`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	} {
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
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

