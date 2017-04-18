package eval

import (
	"bytes"
	"fmt"
	"hash/fnv"
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
	if method == "methods" {
		return h.methods()
	}
	m, ok := hashMethods[method]
	if !ok {
		return newError(NOMETHODERROR, method, h.Type())
	}
	return m(h, args...)
}

func (h *Hash) methods() Object {
	methods := &Array{}
	for key, _ := range hashMethods {
		m := &String{Value: key}
		methods.Members = append(methods.Members, m)
	}
	return methods
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

var hashMethods = map[string]func(h *Hash, args ...Object) Object{
	"keys": func(h *Hash, args ...Object) Object {
		keys := &Array{}
		for _, pair := range h.Pairs {
			keys.Members = append(keys.Members, pair.Key)
		}
		return keys
	},
	"values": func(h *Hash, args ...Object) Object {
		values := &Array{}
		for _, pair := range h.Pairs {
			values.Members = append(values.Members, pair.Value)
		}
		return values
	},
	"push": func(h *Hash, args ...Object) Object {
		if len(args) != 2 {
			return newError(ARGUMENTERROR, "2", len(args))
		}
		if hashable, ok := args[0].(Hashable); ok {
			h.Pairs[hashable.HashKey()] = HashPair{Key: args[0], Value: args[1]}
		} else {
			return newError(KEYERROR, args[0].Type())
		}
		return h
	},
	"pop": func(h *Hash, args ...Object) Object {
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
	},
}
