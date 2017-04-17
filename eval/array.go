package eval

import (
	"bytes"
	"monkey/ast"
	"strings"
)

type Array struct {
	Members []Object
}

func (a *Array) Inspect() string {
	var out bytes.Buffer
	members := []string{}
	for _, m := range a.Members {
		members = append(members, m.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(members, ", "))
	out.WriteString("]")

	return out.String()
}
func (a *Array) Type() ObjectType { return ARRAY_OBJ }

func (a *Array) CallMethod(method string, args []Object) Object {
	if method == "methods" {
		return a.methods()
	}
	m, ok := arrayMethods[method]
	if !ok {
		return newError(NOMETHODERROR, method, a.Type())
	}
	return m(a, args...)
}

func (a *Array) methods() Object {
	methods := &Array{}
	for key, _ := range arrayMethods {
		m := &String{Value: key}
		methods.Members = append(methods.Members, m)
	}
	return methods
}

var arrayMethods = map[string]func(a *Array, args ...Object) Object{
	"pop": func(a *Array, args ...Object) Object {
		l := len(args)
		if l == 0 {
			last := len(a.Members) - 1
			if last < 0 {
				return newError(INDEXERROR, last)
			}
			popped := a.Members[last]
			a.Members = a.Members[:last]
			return popped
		}
		idx := args[1].(*Integer).Value
		if idx == 1 {
			popped, shifted := a.Members[0], a.Members[1:]
			a.Members = shifted
			return popped
		}
		popped := a.Members[idx]
		a.Members = append(a.Members[:idx], a.Members[idx+1:]...)
		return popped
	},
	"push": func(a *Array, args ...Object) Object {
		l := len(args)
		if l != 1 {
			return newError(ARGUMENTERROR, "1", l)
		}
		a.Members = append(a.Members, args[0])
		return a
	},
	"map": func(a *Array, args ...Object) Object {
		if len(args) != 1 {
			return newError(ARGUMENTERROR, "1", len(args))
		}
		block, ok := args[0].(*Function)
		if !ok {
			return newError(INPUTERROR, args[0].Type(), "map")
		}
		arr := &Array{}
		arr.Members = []Object{}
		s := NewScope(nil)
		for _, argument := range a.Members {
			s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument)
			r := Eval(block.Literal.Body, s)
			if obj, ok := r.(*ReturnValue); ok {
				r = obj.Value
			}
			arr.Members = append(arr.Members, r)
		}
		return arr
	},
	"filter": func(a *Array, args ...Object) Object {
		if len(args) != 1 {
			return newError(ARGUMENTERROR, "1", len(args))
		}
		block, ok := args[0].(*Function)
		if !ok {
			return newError(INPUTERROR, args[0].Type(), "filter")
		}
		arr := &Array{}
		arr.Members = []Object{}
		s := NewScope(nil)
		for _, argument := range a.Members {
			s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument)
			r, ok := Eval(block.Literal.Body, s).(*Boolean)
			if !ok {
				return newError(RTERROR, "BOOLEAN")
			}
			if r.Value {
				arr.Members = append(arr.Members, argument)
			}
		}
		return arr
	},
}
