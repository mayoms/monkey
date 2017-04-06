package eval

import (
	"fmt"
)

var builtins map[string]*Builtin

func init() {
	builtins = map[string]*Builtin{
		"len": &Builtin{
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
		"pop": &Builtin{
			Fn: func(args ...Object) Object {
				l := len(args)
				if !(l == 1 || l == 2) {
					return &Error{Message: fmt.Sprintf("too many arguments. expected=1 or 2. got=%d", len(args))}
				}
				switch obj := args[0].(type) {
				case *Array:
					if l == 1 {
						last := len(obj.Members) - 1
						popped := obj.Members[last]
						obj.Members = obj.Members[:last]
						return popped
					}
					idx := args[1].(*Integer).Value
					if idx == 0 {
						popped, shifted := obj.Members[0], obj.Members[1:]
						obj.Members = shifted
						return popped
					}
					popped := obj.Members[idx]
					obj.Members = append(obj.Members[:idx], obj.Members[idx+1:]...)
					return popped

				default:

					return &Error{Message: fmt.Sprintf("unsupported type %T", args[0])}
				}
			},
		},
		"push": &Builtin{
			Fn: func(args ...Object) Object {
				l := len(args)
				if l != 2 {
					return &Error{Message: fmt.Sprintf("too many arguments. expected=1. got=%d", len(args))}
				}
				switch obj := args[0].(type) {
				case *Array:
					obj.Members = append(obj.Members, args[1])
					return obj
				default:
					return &Error{Message: fmt.Sprintf("unsupported type %T", args[1])}
				}
			},
		},
		"map": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return &Error{Message: fmt.Sprintf("too many arguments. expected=2. got=%d", len(args))}
				}
				array, ok := args[0].(*Array)
				if !ok {
					return &Error{Message: fmt.Sprintf("unsupported type %T", args[0])}
				}
				block, ok := args[1].(*Function)
				if !ok {
					return &Error{Message: fmt.Sprintf("unsupported input type %T", args[1])}
				}
				a := &Array{}
				a.Members = []Object{}
				s := NewScope(nil)
				for _, argument := range array.Members {
					s.Set(block.Literal.Parameters[0].Value, argument)
					r := Eval(block.Literal.Body, s)
					if obj, ok := r.(*ReturnValue); ok {
						r = obj.Value
					}
					a.Members = append(a.Members, r)
				}
				return a
			},
		},
	}
}
