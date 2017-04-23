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
	INCLUDED_OBJ     = "INCLUDE"
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

type IncludedObject struct {
	Name  string
	Scope *Scope
}

func (io *IncludedObject) Inspect() string  { return fmt.Sprintf("included object: %s", io.Name) }
func (io *IncludedObject) Type() ObjectType { return INCLUDED_OBJ }
func (io *IncludedObject) CallMethod(method string, args []Object) Object {
	return newError(NOMETHODERROR, method, io.Type())
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
