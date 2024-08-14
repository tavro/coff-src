package eval

import (
	"coff-src/src/coff/object"
	"coff-src/src/coff/ast"
)

var (
	NULL = &object.Null{}
	TRUE = &object.Bool{Value: true}
	FALSE = &object.Bool{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntLiteral:
		return &object.Int{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBoolObject(node.Value)
	}
	
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func nativeBoolToBoolObject(input bool) *object.Bool {
	if input {
		return TRUE
	}

	return FALSE
}