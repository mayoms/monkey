package eval

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type BuiltinFunc func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunc
}

var builtins map[string]*Builtin

func init() {
	builtins = map[string]*Builtin{
		"abs": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				i, ok := args[0].(*Integer)
				if !ok {
					return newError(INPUTERROR, args[0].Type(), "abs")
				}
				if i.Value > -1 {
					return i
				}
				return &Integer{Value: i.Value * -1}
			},
		},
		"addm": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 3 {
					return newError(ARGUMENTERROR, "2", len(args))
				}
				st, ok := args[0].(*Struct)
				if !ok {
					return newError(CONSTRUCTERR, "first", st.Type(), args[0].Type())
				}
				name, ok := args[1].(*String)
				if !ok {
					return newError(CONSTRUCTERR, "second", name.Type(), args[1].Type())
				}
				fn, ok := args[2].(*Function)
				if !ok {
					return newError(CONSTRUCTERR, "second", name.Type(), args[1].Type())
				}
				st.methods[name.Value] = fn
				return NULL
			},
		},
		"chr": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				i, ok := args[0].(*Integer)
				if !ok {
					return newError(INPUTERROR, args[0].Type(), "chr")
				}
				if i.Value < 0 || i.Value > 255 {
					return newError(INPUTERROR, i.Inspect(), "chr")
				}
				return &String{Value: string(i.Value)}
			},
		},
		"open": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				s, ok := args[0].(*String)
				if !ok {
					return newError(INPUTERROR, args[0].Type(), "ord")
				}
				f, err := os.Open(s.Value)
				if err != nil {
					return &Error{Message: err.Error()}
				}
				return &FileObject{File: f}
			},
		},
		"int": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				switch input := args[0].(type) {
				case *Integer:
					return input
				case *String:
					n, err := strconv.Atoi(input.Value)
					if err != nil {
						return newError(INPUTERROR, "STRING: "+input.Value, "int")
					}
					return &Integer{Value: int64(n)}
				}
				return newError(INPUTERROR, args[0].Type(), "int")
			},
		},
		"str": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				switch input := args[0].(type) {
				case *String:
					return input
				default:
					return &String{Value: input.Inspect()}
				}
				return newError(INPUTERROR, args[0].Type(), "str")
			},
		},
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
		"methods": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				methods := &Array{}
				t := reflect.TypeOf(args[0])
				for i := 0; i < t.NumMethod(); i++ {
					m := t.Method(i).Name
					if !(m == "Type" || m == "CallMethod" || m == "HashKey" || m == "Inspect") {
						methods.Members = append(methods.Members, &String{Value: strings.ToLower(m)})
					}
				}
				return methods
			},
		},
		"ord": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				s, ok := args[0].(*String)
				if !ok {
					return newError(INPUTERROR, args[0].Type(), "ord")
				}
				if len(s.Value) > 1 {
					return newError(INLENERR, "1", len(s.Value))
				}
				return &Integer{Value: int64(s.Value[0])}
			},
		},
		"puts": &Builtin{
			Fn: func(args ...Object) Object {
				fmt.Println(args[0].Inspect())
				return NULL
			},
		},
		"type": &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError(ARGUMENTERROR, "1", len(args))
				}
				return &String{Value: fmt.Sprintf("%s", args[0].Type())}
			},
		},
	}
}
