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
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.StrLiteral:
		return &object.Str{Value: node.Value}
	case *ast.ArrLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Arr{Elements: elements}
	case *ast.IdxExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIdxExpression(left, index)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	}
	
	return nil
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Env,) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
	
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
	
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIdxExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARR_OBJ && index.Type() == object.INT_OBJ:
		return evalArrIdxExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIdxExpression(left, index)
	default:
		return newError("index operator is not supported: %s", left.Type())
	}
}

func evalHashIdxExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	
	return pair.Value
}

func evalArrIdxExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Arr)
	
	idx := index.(*object.Int).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	
	return arrayObject.Elements[idx]
}

func applyFunction(fun object.Object, args []object.Object) object.Object {
	switch fun := fun.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fun, args)
		evaluated := Eval(fun.Body, extendedEnv)
		return unwrapRetVal(evaluated)
	case *object.Std:
		return fun.Fun(args...)
	default:
		return newError("not a function: %s", fun.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object,) *object.Env {
	env := object.NewEnclosedEnv(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	
	return env
}

func unwrapRetVal(obj object.Object) object.Object {
	if retVal, ok := obj.(*object.RetVal); ok {
		return retVal.Value
	}
	
	return obj
}

func evalExpressions(exps []ast.Expression, env *object.Env,) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Env,) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	
	if std, ok := stds[node.Value]; ok {
		return std
	}
	
	return newError("identifier is not found: " + node.Value)
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
	case left.Type() == object.STR_OBJ && right.Type() == object.STR_OBJ:
		return evalStrInfixExpression(operator, left, right)
	}
}

func evalStrInfixExpression(operator string, left, right object.Object,) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.Str).Value
	rightVal := right.(*object.Str).Value

	return &object.Str{Value: leftVal + rightVal}
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

