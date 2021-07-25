package evaluator

import (
	"github.com/wawoon/monkeylang/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return object.MakeInt(int64(len(arg.Value)))
			case *object.Array:
				return object.MakeInt(int64(len(arg.Elements)))
			case *object.Null:
				return object.MakeInt(int64(0))
			default:
				return newError("argument to `len` not supported, got %s", arg.Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return arg.Elements[0]
			default:
				return newError("argument to `first` not supported, got %s", arg.Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return arg.Elements[len(arg.Elements)-1]
			default:
				return newError("argument to `last` not supported, got %s", arg.Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}

				newElements := make([]object.Object, len(arg.Elements)-1)
				copy(newElements, arg.Elements[1:])
				return &object.Array{Elements: newElements}
			default:
				return newError("argument to `rest` not supported, got %s", arg.Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				newElements := make([]object.Object, len(arg.Elements)+1)
				copy(newElements, arg.Elements)
				newElements[len(arg.Elements)] = args[1]
				return &object.Array{Elements: newElements}
			default:
				return newError("argument to `push` not supported, got %s", arg.Type())
			}
		},
	},
}
