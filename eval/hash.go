package eval

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkey/ast"
	"strings"
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

type Hashable interface {
	HashKey() HashKey
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s-> %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (h *Hash) CallMethod(method string, args []Object) Object {
	switch method {
	case "filter":
		return h.Filter(args...)
	case "keys":
		return h.Keys(args...)
	case "map":
		return h.Map(args...)
	case "pop":
		return h.Pop(args...)
	case "push":
		return h.Push(args...)
	case "values":
		return h.Values(args...)
	}

	return newError(NOMETHODERROR, method, h.Type())
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (h *Hash) Filter(args ...Object) Object {
	if len(args) != 1 {
		return newError(ARGUMENTERROR, "1", len(args))
	}
	block, ok := args[0].(*Function)
	if !ok {
		return newError(INPUTERROR, args[0].Type(), "filter")
	}
	if len(block.Literal.Parameters) != 2 {
		return newError(ARGUMENTERROR, "2", len(block.Literal.Parameters))
	}
	hash := &Hash{Pairs: make(map[HashKey]HashPair)}
	s := NewScope(nil)
	for _, argument := range h.Pairs {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument.Key)
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, argument.Value)
		r, ok := Eval(block.Literal.Body, s).(*Boolean)
		if !ok {
			return newError(RTERROR, "BOOLEAN")
		}
		if r.Value {
			hash.Push(argument.Key, argument.Value)
		}
	}
	return hash

}

func (h *Hash) Keys(args ...Object) Object {
	keys := &Array{}
	for _, pair := range h.Pairs {
		keys.Members = append(keys.Members, pair.Key)
	}
	return keys
}

func (h *Hash) Map(args ...Object) Object {
	if len(args) != 1 {
		return newError(ARGUMENTERROR, "1", len(args))
	}
	block, ok := args[0].(*Function)
	if !ok {
		return newError(INPUTERROR, args[0].Type(), "map")
	}
	if len(block.Literal.Parameters) != 2 {
		return newError(ARGUMENTERROR, "2", len(block.Literal.Parameters))
	}
	hash := &Hash{Pairs: make(map[HashKey]HashPair)}
	s := NewScope(nil)
	for _, argument := range h.Pairs {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument.Key)
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, argument.Value)
		r := Eval(block.Literal.Body, s)
		if obj, ok := r.(*ReturnValue); ok {
			r = obj.Value
		}
		rh, ok := r.(*Hash)
		if !ok {
			newError(RTERROR, HASH_OBJ)
		}
		for _, v := range rh.Pairs {
			hash.Push(v.Key, v.Value)
		}
	}
	return hash
}

func (h *Hash) Pop(args ...Object) Object {
	if len(args) != 1 {
		return newError(ARGUMENTERROR, "1", len(args))
	}
	hashable, ok := args[0].(Hashable)
	if !ok {
		return newError(KEYERROR, args[0].Type())
	}
	if hashPair, ok := h.Pairs[hashable.HashKey()]; ok {
		delete(h.Pairs, hashable.HashKey())
		return hashPair.Value
	}
	return NULL
}

func (h *Hash) Push(args ...Object) Object {
	if len(args) != 2 {
		return newError(ARGUMENTERROR, "2", len(args))
	}
	if hashable, ok := args[0].(Hashable); ok {
		h.Pairs[hashable.HashKey()] = HashPair{Key: args[0], Value: args[1]}
	} else {
		return newError(KEYERROR, args[0].Type())
	}
	return h
}

func (h *Hash) Values(args ...Object) Object {
	values := &Array{}
	for _, pair := range h.Pairs {
		values.Members = append(values.Members, pair.Value)
	}
	return values
}
