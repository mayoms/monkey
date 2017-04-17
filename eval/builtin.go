package eval

import "fmt"

type BuiltinFunc func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunc
}

var builtins map[string]*Builtin

func init() {
	builtins = map[string]*Builtin{
		"len": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Members))}

				}
				return newError(NOMETHODERROR, "len", args[0].Type())
			},
		},
		"puts": &Builtin{
			Fn: func(args ...Object) Object {
				fmt.Println(args[0].Inspect())
				return NULL
			},
		},
	}
}
