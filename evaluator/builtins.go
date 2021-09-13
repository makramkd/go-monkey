package evaluator

import "github.com/makramkd/go-monkey/object"

var builtins = map[string]*object.Builtin{
	"len": {
		F: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=%d", len(args), 1)
			}

			a := args[0]

			switch a.Type() {
			case object.STRING:
				return &object.Integer{Value: int64(len(a.(*object.String).Value))}
			case object.ARRAY:
				return &object.Integer{Value: int64(len(a.(*object.Array).Values))}
			default:
				return newError("argument to 'len' not supported, got %s", a.Type())
			}
		},
	},
	"first": {
		F: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=%d", len(args), 1)
			}

			a := args[0]

			switch a.Type() {
			case object.ARRAY:
				arr := a.(*object.Array)
				if len(arr.Values) > 0 {
					return arr.Values[0]
				}
				return NULL
			default:
				return newError("argument to 'first' must be ARRAY, got %s", a.Type())
			}
		},
	},
	"last": {
		F: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=%d", len(args), 1)
			}

			a := args[0]

			switch a.Type() {
			case object.ARRAY:
				arr := a.(*object.Array)
				if len(arr.Values) > 0 {
					return arr.Values[len(arr.Values)-1]
				}
				return NULL
			default:
				return newError("argument to 'first' must be ARRAY, got %s", a.Type())
			}
		},
	},
	"rest": {
		F: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=%d", len(args), 1)
			}

			a := args[0]

			switch a.Type() {
			case object.ARRAY:
				arr := a.(*object.Array)
				length := len(arr.Values)
				if length > 0 {
					ret := make([]object.Object, length-1)
					copy(ret, arr.Values[1:length])
					return &object.Array{Values: ret}
				}
				return NULL
			default:
				return newError("argument to 'first' must be ARRAY, got %s", a.Type())
			}
		},
	},
	"push": {
		F: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=%d", len(args), 1)
			}

			a := args[0]

			switch a.Type() {
			case object.ARRAY:
				arr := a.(*object.Array)
				length := len(arr.Values)

				ret := make([]object.Object, length+1)
				copy(ret, arr.Values)
				ret[len(ret)-1] = args[1]

				return &object.Array{Values: ret}
			default:
				return newError("argument to 'first' must be ARRAY, got %s", a.Type())
			}
		},
	},
}
