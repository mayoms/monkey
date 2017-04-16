package eval

import (
	"fmt"
	"monkey/ast"
)

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
		"pop": &Builtin{
			Fn: func(args ...Object) Object {
				l := len(args)
				if !(l == 1 || l == 2) {
					return newError(ARGUMENTERROR, "1 or 2", len(args))
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

					return newError(NOMETHODERROR, "pop", args[0].Type())
				}
			},
		},
		"push": &Builtin{
			Fn: func(args ...Object) Object {
				l := len(args)
				if l != 2 {
					return newError(ARGUMENTERROR, "1", l-1)
				}
				switch obj := args[0].(type) {
				case *Array:
					obj.Members = append(obj.Members, args[1])
					return obj
				default:
					return newError(NOMETHODERROR, "push", args[0].Type())
				}
			},
		},
		"map": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError(ARGUMENTERROR, "1", len(args)-1)
				}
				array, ok := args[0].(*Array)
				if !ok {
					return newError(NOMETHODERROR, "map", args[0].Type())
				}
				block, ok := args[1].(*Function)
				if !ok {
					return newError(INPUTERROR, args[1].Type(), "map")
				}
				a := &Array{}
				a.Members = []Object{}
				s := NewScope(nil)
				for _, argument := range array.Members {
					s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument)
					r := Eval(block.Literal.Body, s)
					if obj, ok := r.(*ReturnValue); ok {
						r = obj.Value
					}
					a.Members = append(a.Members, r)
				}
				return a
			},
		},
		"filter": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError(ARGUMENTERROR, "1", len(args)-1)
				}
				array, ok := args[0].(*Array)
				if !ok {
					return newError(NOMETHODERROR, "filter", args[0].Type())
				}
				block, ok := args[1].(*Function)
				if !ok {
					return newError(INPUTERROR, args[1].Type(), "filter")
				}
				a := &Array{}
				a.Members = []Object{}
				s := NewScope(nil)
				for _, argument := range array.Members {
					s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument)
					r, ok := Eval(block.Literal.Body, s).(*Boolean)
					if !ok {
						return newError(RTERROR, "BOOLEAN")
					}
					if r.Value {
						a.Members = append(a.Members, argument)
					}
				}
				return a
			},
		},
	}
}
