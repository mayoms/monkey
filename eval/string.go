package eval

import (
	"bytes"
)

type String struct{ Value string }

func (s *String) Inspect() string  { return `"` + s.Value + `"` }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) CallMethod(method string, args []Object) Object {

	switch method {
	case "find":
		return s.find(args...)
	case "lower":
		return s.lower(args...)
	case "reverse":
		return s.reverse(args...)
	case "upper":
		return s.upper(args...)
	case "lstrip":
		return s.lstrip(args...)
	}
	return newError(NOMETHODERROR, method, s.Type())
}

func (s *String) find(args ...Object) Object {
	if len(args) != 1 {
		return newError(ARGUMENTERROR, "1", len(args))
	}
	sub, ok := args[0].(*String)
	if !ok {
		return newError(INPUTERROR, args[0].Type(), "find")
	}
	subl := len(sub.Value)
	strl := len(s.Value)
	if subl == 0 || strl == 0 {
		return NULL
	}
	if subl > strl {
		return NULL
	}
	if s.Value == sub.Value {
		return &Integer{Value: 0}
	}
	for i := range s.Value {
		if s.Value[i:i+subl] == sub.Value {
			return &Integer{Value: int64(i)}
		}
		if i+subl > strl {
			break
		}
	}
	return NULL
}

func (s *String) lower(args ...Object) Object {
	if len(args) != 0 {
		return newError(ARGUMENTERROR, "0", len(args))
	}
	if s.Value == "" {
		return s
	}
	var out bytes.Buffer
	for _, ch := range s.Value {
		if 'A' <= ch && ch <= 'Z' {
			out.WriteRune(ch + 32)
			continue
		}
		out.WriteRune(ch)
	}
	return &String{Value: out.String()}
}

func (s *String) lstrip(args ...Object) Object {
	if len(args) > 1 {
		return newError(ARGUMENTERROR, "1", len(args))
	}
	var substr string
	var subl int

	if len(args) == 1 {
		substr = args[0].(*String).Value
		subl = len(substr)
	}
	sl := len(s.Value)
	if subl > sl {
		return s
	}
	c := 0
	if substr == "" {
		for isWhiteSpace(s.Value[c]) {
			c++
			if c == sl-1 {
				return &String{Value: ""}
			}
		}
		return &String{Value: s.Value[c:]}
	}
	for i := range s.Value {
		if s.Value[i] != substr[i%subl] {
			c = i
			break
		}
	}
	if c == 0 {
		return s
	}
	return &String{Value: s.Value[c:]}
}
func isWhiteSpace(b byte) bool {
	return (b == 0x20 || b == 0x9 || b == 0xa || b == 0xd)
}

func (s *String) reverse(args ...Object) Object {
	if len(args) != 0 {
		return newError(ARGUMENTERROR, "0", len(args))
	}

	end := len(s.Value) - 1
	if end < 1 {
		return s
	}
	var out bytes.Buffer
	for i := range s.Value {
		out.WriteByte(s.Value[end-i])
	}
	return &String{Value: out.String()}
}

// func(s *String) strip(...Object) Object{
// 		kind := "both"
// 		var start int
// 		var end int
// 		if len(args) > 2 {
// 			return newError(ARGUMENTERROR, "1 or 2", len(args))
// 		}
// 		substr, ok := args[0].(*String); !ok {
// 			return newError(INPUTERROR, args[0].Type(), "len")
// 		}
// 		if len(args) == 2 {
// 			k, ok := args[1].(*String); !ok {
// 				return newError(INPUTERROR, args[0].Type(), "strip")
// 				if k.Value != "left" || k.Value != "right" {
// 					return newError(INPUTERROR, "STRING: " + args[0].Value, "strip")
// 				}

//  			}
// 		}
// 	}

func (s *String) upper(args ...Object) Object {
	if len(args) != 0 {
		return newError(ARGUMENTERROR, "0", len(args))
	}
	if s.Value == "" {
		return s
	}
	var out bytes.Buffer
	for _, ch := range s.Value {
		if 'a' <= ch && ch <= 'z' {
			out.WriteRune(ch - 32)
			continue
		}
		out.WriteRune(ch)
	}
	return &String{Value: out.String()}
}
