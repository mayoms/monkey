package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY_OBJ"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type BuiltinFunc func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunc
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

type Function struct {
	Literal *ast.FunctionLiteral
	Scope   *Scope
}

func (f *Function) Inspect() string  { return f.Literal.String() }
func (r *Function) Type() ObjectType { return FUNCTION_OBJ }

type Array struct {
	Value   string
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
	out.WriteString("\n")

	return out.String()
}
func (a *Array) Type() ObjectType { return ARRAY_OBJ }

type ReturnValue struct{ Value Object }

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Integer struct{ Value int64 }

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type String struct{ Value string }

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }

type Boolean struct{ Value bool }

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%v", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type Error struct{ Message string }

func (e *Error) Inspect() string  { return "Err: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }
