package eval

import (
	"coff-src/src/coff/object"
)

var stds = map[string]*object.Std{
	"len": &object.Std{
		Fun: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			
			switch arg := args[0].(type) {
			case *object.Str:
				return &object.Int{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` is not supported, got %s", args[0].Type())
			}
		},
	},
}