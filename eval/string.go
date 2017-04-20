package eval

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
}
