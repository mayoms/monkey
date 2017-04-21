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
		return s.Find(args...)
	case "lower":
		return s.Lower(args...)
	case "reverse":
		return s.Reverse(args...)
	case "upper":
		return s.Upper(args...)
	case "lstrip":
		return s.Lstrip(args...)
	case "rstrip":
		return s.Rstrip(args...)
	case "strip":
		return s.Strip(args...)
	case "split":
		return s.Split(args...)
	}
	return newError(NOMETHODERROR, method, s.Type())
}

func (s *String) Find(args ...Object) Object {
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

func (s *String) Lower(args ...Object) Object {
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

func (s *String) Lstrip(args ...Object) Object {
	if len(args) > 1 {
		return newError(ARGUMENTERROR, "1", len(args))
	}
	var substr string
	var subl int

	if len(args) == 1 {
		sObj, ok := args[0].(*String)
		if !ok {
			return newError(INPUTERROR, args[0].Type(), "rstrip")
		}
		substr = sObj.Value
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

func (s *String) Reverse(args ...Object) Object {
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

func (s *String) Rstrip(args ...Object) Object {
	if len(args) > 1 {
		return newError(ARGUMENTERROR, "1", len(args))
	}
	var substr string
	var subl int

	if len(args) == 1 {
		sObj, ok := args[0].(*String)
		if !ok {
			return newError(INPUTERROR, args[0].Type(), "rstrip")
		}
		substr = sObj.Value
		subl = len(substr)
	}
	sl := len(s.Value) - 1
	if subl-1 > sl {
		return s
	}
	if substr == "" {
		for isWhiteSpace(s.Value[sl]) {
			sl--
			if sl == 0 {
				return &String{Value: ""}
			}
		}
		return &String{Value: s.Value[:sl+1]}
	}
	for i := range s.Value {
		i = sl - 1 - i
		if s.Value[i] != substr[i%subl] {
			sl = i
			break
		}
	}
	if sl == 0 {
		return s
	}
	return &String{Value: s.Value[:sl+1]}
}

func (s *String) Strip(args ...Object) Object {
	l := s.Lstrip(args...)
	return l.(*String).Rstrip(args...)
}

func (s *String) Split(args ...Object) Object {
	if len(args) > 1 {
		return newError(ARGUMENTERROR, "0 or 1", len(args))
	}

	var del string

	if len(args) == 1 {
		sObj, ok := args[0].(*String)
		if !ok {
			return newError(INPUTERROR, args[0].Type(), "rstrip")
		}
		del = sObj.Value
	}
	dl := len(del)
	sl := len(s.Value)

	if dl > sl {
		return s
	}

	if del == "" {
		return splitWhiteSpace(s.Value)
	}

	start := 0
	stop := sl
	a := &Array{}
	for i := 0; i < sl; i++ {
		if s.Value[i] == del[0] {
			stop = i
			a.Members = append(a.Members, &String{Value: s.Value[start:stop]})
			for p := 0; p < dl; p++ {
				if s.Value[p] != del[p%dl] {
					start = i
				}
				i++
			}
			if start < i {
				start = stop + 1
			}
		}
	}
	a.Members = append(a.Members, &String{Value: s.Value[start:]})
	return a
}

func splitWhiteSpace(s string) Object {
	sl := len(s)
	start := 0
	stop := sl
	a := &Array{}

	for i := 0; i < sl; i++ {
		if isWhiteSpace(s[i]) {
			stop = i
			a.Members = append(a.Members, &String{Value: s[start:stop]})
			for p := i; p < sl; p++ {
				if !isWhiteSpace(s[p]) {
					start = p
					break
				}
				i++
			}
			if start < i {
				start = stop + 1
			}
		}
	}
	a.Members = append(a.Members, &String{Value: s[start:]})
	return a
}

func (s *String) Upper(args ...Object) Object {
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
