package eval

import (
	"fmt"
	"monkey/ast"
)

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	NULL  = &Null{}
)

func Eval(node ast.Node, scope *Scope) Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, scope)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, scope)
	case *ast.LetStatement:
		return evalLetStatement(node, scope)
	case *ast.ReturnStatement:
		return evalReturnStatment(node, scope)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	case *ast.StringLiteral:
		return evalStringLiteral(node)
	case *ast.Identifier:
		return evalIdentifier(node, scope)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, scope)
	case *ast.HashLiteral:
		return evalHashLiteral(node, scope)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, scope)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, scope)
	case *ast.InfixExpression:
		return evalInfixExpression(node, scope)
	case *ast.IfExpression:
		return evalIfExpression(node, scope)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, scope)
	case *ast.CallExpression:
		return evalFunctionCall(node, scope)
	case *ast.MethodCallExpression:
		return evalMethodCallExpression(node, scope)
	case *ast.IndexExpression:
		return evalIndexExpression(node, scope)
	}
	return nil
}

func evalLetStatement(l *ast.LetStatement, scope *Scope) (val Object) {
	if val = Eval(l.Value, scope); val.Type() != ERROR_OBJ {
		return scope.Set(l.Name.String(), val)
	}
	return
}

func evalReturnStatment(r *ast.ReturnStatement, scope *Scope) Object {
	if value := Eval(r.ReturnValue, scope); value != nil {
		return &ReturnValue{Value: value}
	}
	return NULL
}

func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalIntegerLiteral(i *ast.IntegerLiteral) Object {
	return &Integer{Value: i.Value}
}

func evalStringLiteral(s *ast.StringLiteral) Object {
	return &String{Value: s.Value}
}

func evalArrayLiteral(a *ast.ArrayLiteral, scope *Scope) Object {
	return &Array{Members: evalArgs(a.Members, scope)}
}

func evalIdentifier(i *ast.Identifier, scope *Scope) Object {
	if val, ok := scope.Get(i.String()); ok {
		return val
	}
	return &Error{Message: fmt.Sprintf("unknown identifier: %s", i.String())}
}

func evalHashLiteral(hl *ast.HashLiteral, scope *Scope) Object {
	hashMap := make(map[HashKey]HashPair)
	// TODO: { 1 -> true, 2 -> "five", "three"-> "four } causes some sort of infinite loop
	for key, value := range hl.Pairs {
		key := Eval(key, scope)
		if hashable, ok := key.(Hashable); ok {
			hashMap[hashable.HashKey()] = HashPair{Key: key, Value: Eval(value, scope)}
		} else {
			return &Error{Message: fmt.Sprintf("key error: %T not hashable", key.Type())}
		}
	}
	return &Hash{Pairs: hashMap}
}

func evalFunctionLiteral(fl *ast.FunctionLiteral, scope *Scope) Object {
	return &Function{Literal: fl, Scope: scope}
}

func evalPrefixExpression(p *ast.PrefixExpression, s *Scope) Object {
	right := Eval(p.Right, s)
	if right.Type() == ERROR_OBJ {
		return right
	}
	switch p.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		if i, ok := right.(*Integer); ok {
			i.Value = -i.Value
			return i
		}
	}
	return &Error{Message: fmt.Sprintf("unknown operator: %s%s", p.Operator, right.Type())}
}

func evalInfixExpression(i *ast.InfixExpression, s *Scope) Object {
	left := Eval(i.Left, s)
	right := Eval(i.Right, s)
	if left.Type() == ERROR_OBJ {
		return left
	} else if right.Type() == ERROR_OBJ {
		return right
	}

	var errMsg string
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntInfixExpression(i.Operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(i.Operator, left, right)
	case i.Operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case i.Operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		errMsg = fmt.Sprintf("type mismatch: %s %s %s", left.Type(), i.Operator, right.Type())
	default:
		errMsg = fmt.Sprintf("unknown operator: %s %s %s", left.Type(), i.Operator, right.Type())
	}
	return &Error{Message: errMsg}
}

func evalIfExpression(ie *ast.IfExpression, s *Scope) Object {
	condition := Eval(ie.Condition, s)
	if isTrue(condition) {
		return evalBlockStatements(ie.Consequence.Statements, s)
	} else if ie.Alternative != nil {
		return evalBlockStatements(ie.Alternative.Statements, s)
	}
	return NULL
}

func evalBlockStatements(block []ast.Statement, scope *Scope) (results Object) {
	for _, statement := range block {
		results = Eval(statement, scope)
		if results != nil && results.Type() == RETURN_VALUE_OBJ {
			return
		}
	}
	return
}

func evalFunctionCall(call *ast.CallExpression, s *Scope) Object {
	fn, ok := s.Get(call.Function.String())
	if !ok {
		if f, ok := call.Function.(*ast.FunctionLiteral); ok {
			fn = &Function{Literal: f, Scope: s}
			s.Set(call.Function.String(), fn)
		} else if builtin, ok := builtins[call.Function.String()]; ok {
			return builtin.Fn(evalArgs(call.Arguments, s)...)
		} else {
			return &Error{Message: fmt.Sprintf("unknown identifier: %s", call.Function.String())}
		}
	}
	f := fn.(*Function)
	f.Scope = NewScope(s)
	args := evalArgs(call.Arguments, f.Scope)
	// TODO: If the wrong number of params is passed a panic occurs
	for i, v := range f.Literal.Parameters {
		f.Scope.Set(v.String(), args[i])
	}
	r := Eval(f.Literal.Body, f.Scope)
	if obj, ok := r.(*ReturnValue); ok {
		return obj.Value
	}
	return r
}

func evalMethodCallExpression(call *ast.MethodCallExpression, scope *Scope) Object {
	obj := Eval(call.Object, scope)
	method, ok := call.Call.(*ast.CallExpression)
	if !ok {
		return &Error{Message: fmt.Sprintf("Method call not *ast.CallExpression. got=%T", call.Call)}
	}
	args := evalArgs(method.Arguments, scope)
	return obj.CallMethod(method.Function.String(), args)
}

func evalArgs(args []ast.Expression, scope *Scope) []Object {
	e := []Object{}
	for _, v := range args {
		e = append(e, Eval(v, scope))
	}
	return e
}

func evalProgram(program *ast.Program, scope *Scope) (results Object) {
	for _, statement := range program.Statements {
		results = Eval(statement, scope)
		switch s := results.(type) {
		case *ReturnValue:
			return s.Value
		case *Error:
			return s
		}
	}
	return
}

func evalIndexExpression(ie *ast.IndexExpression, s *Scope) Object {
	left := Eval(ie.Left, s)
	index := Eval(ie.Index, s)
	if index.Type() == ERROR_OBJ {
		return index
	}
	switch iterable := left.(type) {
	case *Array:
		return evalArrayIndex(iterable, index)
	case *Hash:
		return evalHashKeyIndex(iterable, index)
	}
	return &Error{Message: fmt.Sprintf("index not supported for type: %s", left.Type())}
}

func evalHashKeyIndex(hash *Hash, key Object) Object {
	hashable, ok := key.(Hashable)
	if !ok {
		return &Error{Message: fmt.Sprintf("key error: %s not hashable", key.Type())}
	}
	hashPair, ok := hash.Pairs[hashable.HashKey()]
	if !ok {
		return NULL
	}
	return hashPair.Value
}

func evalArrayIndex(array *Array, index Object) Object {
	idx, ok := index.(*Integer)
	if !ok {
		return &Error{Message: fmt.Sprintf("type error: index not integer. got=%s", index.Type())}
	}
	length := int64(len(array.Members))
	if idx.Value > length-1 || idx.Value < -length {
		return &Error{Message: fmt.Sprintf("index %d out of range", idx.Value)}
	}
	if idx.Value < 0 {
		return array.Members[(length)+idx.Value]
	}
	return array.Members[idx.Value]
}

func isTrue(obj Object) bool {
	switch obj {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func evalStringInfixExpression(operator string, left Object, right Object) Object {
	l := left.(*String)
	r := right.(*String)

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(l.Value == r.Value)
	case "!=":
		return nativeBoolToBooleanObject(l.Value != r.Value)
	case "+":
		// TODO: "Hello, + "World" causes some sort of infinite loop
		return &String{Value: l.Value + r.Value}
	}
	return &Error{Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, left.Type())}
}

func evalIntInfixExpression(operator string, left Object, right Object) Object {
	l := left.(*Integer)
	r := right.(*Integer)

	switch operator {
	case "+":
		return &Integer{Value: l.Value + r.Value}
	case "-":
		return &Integer{Value: l.Value - r.Value}
	case "*":
		return &Integer{Value: l.Value * r.Value}
	case "/":
		return &Integer{Value: l.Value / r.Value}
	case ">":
		return nativeBoolToBooleanObject(l.Value > r.Value)
	case "<":
		return nativeBoolToBooleanObject(l.Value < r.Value)
	case "==":
		return nativeBoolToBooleanObject(l.Value == r.Value)
	case "!=":
		return nativeBoolToBooleanObject(l.Value != r.Value)
	}
	return NULL
}

func evalBangOperatorExpression(right Object) Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}
