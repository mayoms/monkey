package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

type ObjectType string

const (
	// INTEGER_OBJ/*_OBJ = object types
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY_OBJ"

	// NOMETHODERROR = no method on object
	NOMETHODERROR = "No method %s for object %s"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	CallMethod(method string, args []Object) Object
}

type BuiltinFunc func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunc
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) CallMethod(method string, args []Object) Object {
	return &Error{Message: fmt.Sprintf(NOMETHODERROR, method, b.Type())}
}

type Function struct {
	Literal *ast.FunctionLiteral
	Scope   *Scope
}

func (f *Function) Inspect() string  { return f.Literal.String() }
func (r *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) CallMethod(method string, args []Object) Object {
	return &Error{Message: fmt.Sprintf(NOMETHODERROR, method, f.Type())}
}

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
	builtin, ok := Builtins[method]
	if !ok {
		return &Error{Message: fmt.Sprintf(NOMETHODERROR, method, a.Type())}
	}
	args = append([]Object{a}, args...)
	return builtin.Fn(args...)
}

type ReturnValue struct{ Value Object }

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) CallMethod(method string, args []Object) Object {
	return &Error{Message: fmt.Sprintf(NOMETHODERROR, method, rv.Type())}
}

type Integer struct{ Value int64 }

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) CallMethod(method string, args []Object) Object {
	return &Error{Message: fmt.Sprintf(NOMETHODERROR, method, i.Type())}
}

type String struct{ Value string }

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) CallMethod(method string, args []Object) Object {
	builtin, ok := Builtins[method]
	if !ok {
		return &Error{Message: fmt.Sprintf(NOMETHODERROR, method, s.Type())}
	}
	args = append([]Object{s}, args...)
	return builtin.Fn(args...)
}

type Boolean struct{ Value bool }

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%v", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) CallMethod(method string, args []Object) Object {
	return &Error{Message: fmt.Sprintf(NOMETHODERROR, method, b.Type())}
}

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) CallMethod(method string, args []Object) Object {
	return &Error{Message: fmt.Sprintf(NOMETHODERROR, method, n.Type())}
}

type Error struct{ Message string }

func (e *Error) Inspect() string  { return "Err: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) CallMethod(method string, args []Object) Object {
	return &Error{Message: fmt.Sprintf(NOMETHODERROR, method, e.Type())}
}
