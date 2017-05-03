package eval

import (
	"bufio"
	"bytes"
	"fmt"
	"monkey/ast"
	"os"
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
	STRUCT_OBJ       = "STRUCT"
	FILE_OBJ         = "FILE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	CallMethod(method string, args ...Object) Object
}

type FileObject struct {
	File *os.File
	Name string
}

func (f *FileObject) Inspect() string  { return f.Name }
func (f *FileObject) Type() ObjectType { return FILE_OBJ }
func (f *FileObject) CallMethod(method string, args ...Object) Object {
	switch method {
	case "close":
		f.File.Close()
		return NULL
	case "read":
		return f.Read(args...)
	default:
		return newError(NOMETHODERROR, method, f.Type())
	}
}

func (f *FileObject) Read(args ...Object) Object {
	if len(args) != 0 {
		return newError(ARGUMENTERROR, "0", len(args))
	}
	fs := bufio.NewScanner(f.File)
	var out bytes.Buffer
	for {
		scanned := fs.Scan()
		out.WriteString(fs.Text())
		if !scanned {
			break
		}
		if err := fs.Err(); err != nil {
			return &Error{Message: err.Error()}
		}
	}
	return &String{Value: out.String()}
}

type Struct struct {
	Scope   *Scope
	methods map[string]*Function
}

func (s *Struct) Inspect() string {
	var out bytes.Buffer
	out.WriteString("( ")
	for k, v := range s.Scope.store {
		out.WriteString(k)
		out.WriteString("->")
		out.WriteString(v.Inspect())
		out.WriteString(" ")
	}
	out.WriteString(" )")

	return out.String()
}

func (s *Struct) Type() ObjectType { return STRUCT_OBJ }
func (s *Struct) CallMethod(method string, args ...Object) Object {
	fn, ok := s.methods[method]
	if !ok {
		return newError(NOMETHODERROR, method, s.Type())
	}
	fn.Scope = NewScope(nil)
	fn.Scope.Set("self", s)
	for i, v := range fn.Literal.Parameters {
		fn.Scope.Set(v.String(), args[i])
	}
	r := Eval(fn.Literal.Body, fn.Scope)
	if obj, ok := r.(*ReturnValue); ok {
		return obj.Value
	}
	return r
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) CallMethod(method string, args ...Object) Object {
	return newError(NOMETHODERROR, method, b.Type())
}

type IncludedObject struct {
	Name  string
	Scope *Scope
}

func (io *IncludedObject) Inspect() string  { return fmt.Sprintf("included object: %s", io.Name) }
func (io *IncludedObject) Type() ObjectType { return INCLUDED_OBJ }
func (io *IncludedObject) CallMethod(method string, args ...Object) Object {
	return newError(NOMETHODERROR, method, io.Type())
}

type Function struct {
	Literal *ast.FunctionLiteral
	Scope   *Scope
}

func (f *Function) Inspect() string  { return f.Literal.String() }
func (r *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) CallMethod(method string, args ...Object) Object {
	return newError(NOMETHODERROR, method, f.Type())
}

type ReturnValue struct{ Value Object }

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) CallMethod(method string, args ...Object) Object {
	return newError(NOMETHODERROR, method, rv.Type())
}

type Integer struct{ Value int64 }

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) CallMethod(method string, args ...Object) Object {
	return newError(NOMETHODERROR, method, i.Type())
}

type Boolean struct{ Value bool }

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%v", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) CallMethod(method string, args ...Object) Object {
	return newError(NOMETHODERROR, method, b.Type())
}

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) CallMethod(method string, args ...Object) Object {
	return newError(NOMETHODERROR, method, n.Type())
}
