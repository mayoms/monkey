package eval

import "fmt"

// constants for error types
const (
	_ int = iota
	PREFIXOP
	INFIXOP
	UNKNOWNIDENT
	NOMETHODERROR
	NOINDEXERROR
	KEYERROR
	INDEXERROR
	SLICEERROR
	ARGUMENTERROR
	INPUTERROR
	RTERROR
)

var errorType = map[int]string{
	PREFIXOP:      "unsupported operator for prefix expression:'%s' and type: %s",
	INFIXOP:       "unsupported operator for infix expression: '%s' and types %s and %s",
	UNKNOWNIDENT:  "unknown identifier: '%s' is not defined",
	NOMETHODERROR: "undefined method '%s' for object %s",
	NOINDEXERROR:  "index error: type %s is not indexable",
	KEYERROR:      "key error: type %s is not hashable",
	INDEXERROR:    "index error: '%d' out of range",
	SLICEERROR:    "index error: slice '%d:%d' out of range",
	ARGUMENTERROR: "wrong number of arguments. expected=%s, got=%d",
	INPUTERROR:    "unsupported input type '%s' for function or method %s",
	RTERROR:       "return type should be %s.",
}

func newError(t int, args ...interface{}) Object {
	return &Error{Message: fmt.Sprintf(errorType[t], args...)}
}

type Error struct{ Message string }

func (e *Error) Inspect() string  { return "Err: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) CallMethod(method string, args []Object) Object {
	return newError(NOMETHODERROR, method, e.Type())
}
