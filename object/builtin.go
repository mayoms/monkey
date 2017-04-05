package object

import "fmt"

var Builtins = map[string]Builtin{
	"len": Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("too many arguments. expected=1 got=%d", len(args))}
			}
			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			case *Array:
				return &Integer{Value: int64(len(arg.Members))}

			}
			return &Error{Message: fmt.Sprintf("unsupported type: %T", args[0])}
		},
	},
}
