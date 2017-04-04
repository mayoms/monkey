package eval

import (
	"fmt"
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("too many arguments. expected=1 got=%d", len(args))}
			}
			if string, ok := args[0].(*object.String); ok {
				return &object.Integer{Value: int64(len(string.Value))}
			}
			return &object.Error{Message: fmt.Sprintf("unsupported type: %T", args[0])}
		},
	},
}
