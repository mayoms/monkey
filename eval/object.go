package eval

import (
	"fmt"
	"monkey/ast"
)

type ObjectType string

// INTEGER_OBJ/*_OBJ = object types
const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	CallMethod(method string, args []Object) Object
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) CallMethod(method string, args []Object) Object {
	return newError(NOMETHODERROR, method, b.Type())
}

type Function struct {
	Literal *ast.FunctionLiteral
	Scope   *Scope
}

func (f *Function) Inspect() string  { return f.Literal.String() }
func (r *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) CallMethod(method string, args []Object) Object {
	return newError(NOMETHODERROR, method, f.Type())
}

type ReturnValue struct{ Value Object }

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) CallMethod(method string, args []Object) Object {
	return newError(NOMETHODERROR, method, rv.Type())
}

type Integer struct{ Value int64 }

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) CallMethod(method string, args []Object) Object {
	return newError(NOMETHODERROR, method, i.Type())
}

type String struct{ Value string }

func (s *String) Inspect() string  { return `"` + s.Value + `"` }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) CallMethod(method string, args []Object) Object {
	builtin, ok := builtins[method]
	if !ok {
		return newError(NOMETHODERROR, method, s.Type())
	}
	args = append([]Object{s}, args...)
	return builtin.Fn(args...)
}

type Boolean struct{ Value bool }

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%v", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) CallMethod(method string, args []Object) Object {
	return newError(NOMETHODERROR, method, b.Type())
}

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) CallMethod(method string, args []Object) Object {
	return newError(NOMETHODERROR, method, n.Type())
}
