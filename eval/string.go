package eval

import "bytes"

type String struct{ Value string }

func (s *String) Inspect() string  { return `"` + s.Value + `"` }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) CallMethod(method string, args []Object) Object {
	m, ok := stringMethods[method]
	if !ok {
		return newError(NOMETHODERROR, method, s.Type())
	}
	return m(s, args...)
}

var stringMethods = map[string]func(s *String, args ...Object) Object{
	"find": func(s *String, args ...Object) Object {
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
	},
	"reverse": func(s *String, args ...Object) Object {
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
	},
	"upper": func(s *String, args ...Object) Object {
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
	},
}
