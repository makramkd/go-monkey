package evaluator

import "github.com/makramkd/go-monkey/object"

var builtins = map[string]*object.Builtin{
	"len": {
		F: func(args ...object.Object) object.Object {
			if len(args) > 1 {
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
}
