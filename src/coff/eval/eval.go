package eval

import (
	"coff-src/src/coff/object"
	"coff-src/src/coff/ast"
	"fmt"
)

var (
	NULL = &object.Null{}
	TRUE = &object.Bool{Value: true}
	FALSE = &object.Bool{Value: false}
)

func Eval(node ast.Node, env *object.Env) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.RetStatement:
		val := Eval(node.RetVal, env)
		if isError(val) {
			return val
		}
		return &object.RetVal{Value: val}
	case *ast.IntLiteral:
		return &object.Int{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBoolObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.DefStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	}
	
	return nil
}

func evalIdentifier(node *ast.Identifier, env *object.Env,) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}
	
	return val
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Env,) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RET_VAL_OBJ || rt == object.ERR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Env,) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, left, right object.Object,) object.Object {
	switch {
	case left.Type() == object.INT_OBJ && right.Type() == object.INT_OBJ:
		return evalIntInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBoolObject(left == right)
	case operator == "!=":
		return nativeBoolToBoolObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntInfixExpression(operator string, left, right object.Object,) object.Object {
	leftVal := left.(*object.Int).Value
	rightVal := right.(*object.Int).Value
	switch operator {
	case "+":
		return &object.Int{Value: leftVal + rightVal}
	case "-":
		return &object.Int{Value: leftVal - rightVal}
	case "*":
		return &object.Int{Value: leftVal * rightVal}
	case "/":
		return &object.Int{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBoolObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBoolObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBoolObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBoolObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalFacOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INT_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Int).Value
	return &object.Int{Value: -value}
}

func evalFacOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalProgram(program *ast.Program, env *object.Env) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.RetVal:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func nativeBoolToBoolObject(input bool) *object.Bool {
	if input {
		return TRUE
	}

	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERR_OBJ
	}
	return false
}

