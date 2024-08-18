package eval

import (
	"fmt"
	"coff-src/src/coff/object"
)

var stds = map[string]*object.Std{
	"len": &object.Std{
		Fun: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			
			switch arg := args[0].(type) {
			case *object.Arr:
				return &object.Int{Value: int64(len(arg.Elements))}
			case *object.Str:
				return &object.Int{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` is not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.Std{
		Fun: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARR_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Arr)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": &object.Std{
		Fun: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
		
			if args[0].Type() != object.ARR_OBJ {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}
		
			arr := args[0].(*object.Arr)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"rest": &object.Std{
		Fun: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARR_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Arr)
			length := len(arr.Elements)

			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Arr{Elements: newElements}
			}

			return NULL
		},
	},
	"push": &object.Std{
		Fun: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
		
			if args[0].Type() != object.ARR_OBJ {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
		
			arr := args[0].(*object.Arr)
			length := len(arr.Elements)
			
			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			
			return &object.Arr{Elements: newElements}
		},
	},
	"print": &object.Std{
		Fun: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			
			return NULL
		},
	},
}