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

func (a *Array) CallMethod(method string, args ...Object) Object {
	switch method {
	case "count":
		return a.Count(args...)
	case "filter":
		return a.Filter(args...)
	case "index":
		return a.Index(args...)
	case "map":
		return a.Map(args...)
	case "merge":
		return a.Merge(args...)
	case "push":
		return a.Push(args...)
	case "pop":
		return a.Pop(args...)
	case "reduce":
		return a.Reduce(args...)
	}
	return newError(NOMETHODERROR, a.Type(), method)
}

func (a *Array) Count(args ...Object) Object {
	if len(args) < 1 || len(args) > 1 {
		return newError(ARGUMENTERROR, "1", len(args))
	}
	count := 0
	for _, v := range a.Members {
		switch c := args[0].(type) {
		case *Integer:
			if c.Value == v.(*Integer).Value {
				count++
			}
		case *String:
			if c.Value == v.(*String).Value {
				count++
			}
		default:
			if c == v {
				count++
			}
		}
	}
	return &Integer{Value: int64(count)}
}

func (a *Array) Filter(args ...Object) Object {
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
}

func (a *Array) Index(args ...Object) Object {
	if len(args) < 1 || len(args) > 1 {
		return newError(ARGUMENTERROR, "1", len(args))
	}
	for i, v := range a.Members {
		switch c := args[0].(type) {
		case *Integer:
			if c.Value == v.(*Integer).Value {
				return &Integer{Value: int64(i)}
			}
		case *String:
			if c.Value == v.(*String).Value {
				return &Integer{Value: int64(i)}
			}
		default:
			if c == v {
				return &Integer{Value: int64(i)}
			}
		}
	}
	return NULL
}
func (a *Array) Map(args ...Object) Object {
	if len(args) != 1 {
		return newError(ARGUMENTERROR, "1", len(args))
	}
	block, ok := args[0].(*Function)
	if !ok {
		return newError(INPUTERROR, args[0].Type(), "map")
	}
	arr := &Array{}
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
}

func (a *Array) Merge(args ...Object) Object {
	if len(args) != 1 {
		return newError(ARGUMENTERROR, "1", len(args))
	}
	m, ok := args[0].(*Array)
	if !ok {
		return newError(INPUTERROR, args[0].Type(), "array.merge")
	}
	arr := &Array{}
	for _, v := range a.Members {
		arr.Members = append(arr.Members, v)
	}
	for _, v := range m.Members {
		arr.Members = append(arr.Members, v)
	}
	return arr
}

func (a *Array) Pop(args ...Object) Object {
	last := len(a.Members) - 1
	if len(args) == 0 {
		if last < 0 {
			return newError(INDEXERROR, last)
		}
		popped := a.Members[last]
		a.Members = a.Members[:last]
		return popped
	}
	idx := args[0].(*Integer).Value
	if idx < 0 {
		idx = idx + int64(last+1)
	}
	if idx < 0 || idx > int64(last) {
		return newError(INDEXERROR, idx)
	}
	popped := a.Members[idx]
	a.Members = append(a.Members[:idx], a.Members[idx+1:]...)
	return popped
}

func (a *Array) Push(args ...Object) Object {
	l := len(args)
	if l != 1 {
		return newError(ARGUMENTERROR, "1", l)
	}
	a.Members = append(a.Members, args[0])
	return a
}

func (a *Array) Reduce(args ...Object) Object {
	l := len(args)
	if 1 > 2 || l < 1 {
		return newError(ARGUMENTERROR, "1 or 2", l)
	}

	block, ok := args[0].(*Function)
	if !ok {
		return newError(INPUTERROR, args[0].Type(), "map")
	}
	s := NewScope(nil)
	start := 1
	if l == 1 {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, a.Members[0])
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, a.Members[1])
		start += 1
	} else {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, args[1])
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, a.Members[0])
	}
	r := Eval(block.Literal.Body, s)
	for i := start; i < len(a.Members); i++ {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, r)
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, a.Members[i])
		r = Eval(block.Literal.Body, s)
		if obj, ok := r.(*ReturnValue); ok {
			r = obj.Value
		}
	}
	return r

}
