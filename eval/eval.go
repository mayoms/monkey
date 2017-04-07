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
	case *ast.MethodCallExpression:
		return evalMethodCallExpression(node, scope)
	case *ast.CallExpression:
		f_scope := NewScope(scope)
		fn, ok := scope.Get(node.Function.String())
		if !ok {
			if f, ok := node.Function.(*ast.FunctionLiteral); ok {
				fn = &Function{Literal: f, Scope: scope}
				scope.Set(node.Function.String(), fn)
			}
			if builtin, ok := builtins[node.Function.String()]; ok {
				return builtin.Fn(evalArgs(node.Arguments, scope)...)
			}
		}
		return evalFunctionCall(node, f_scope)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, scope)
	case *ast.HashLiteral:
		return evalHashLiteral(node, scope)
	case *ast.FunctionLiteral:
		return &Function{Literal: node, Scope: scope}
	case *ast.LetStatement:
		val := Eval(node.Value, scope)
		if val.Type() == ERROR_OBJ {
			return val
		}
		return scope.Set(node.Name.String(), val)
	case *ast.Identifier:
		if val, ok := scope.Get(node.String()); ok {
			return val
		}
		return &Error{Message: fmt.Sprintf("unknown identifier: %s", node.String())}
	case *ast.ExpressionStatement:
		return Eval(node.Expression, scope)
	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue, scope)
		if value != nil {
			return &ReturnValue{Value: value}
		}
		return NULL
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, scope)
	case *ast.InfixExpression:
		left := Eval(node.Left, scope)
		right := Eval(node.Right, scope)
		if left.Type() == ERROR_OBJ {
			return left
		} else if right.Type() == ERROR_OBJ {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.PrefixExpression:
		right := Eval(node.Right, scope)
		if right.Type() == ERROR_OBJ {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.IfExpression:
		condition := Eval(node.Condition, scope)
		if isTrue(condition) {
			return evalBlockStatements(node.Consequence.Statements, scope)
		} else if node.Alternative != nil {
			return evalBlockStatements(node.Alternative.Statements, scope)
		}
		return NULL
	case *ast.IntegerLiteral:
		return &Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IndexExpression:
		left := Eval(node.Left, scope)
		index := Eval(node.Index, scope)
		return evalIndexExpression(left, index)
	}
	return nil
}

func evalIndexExpression(left, index Object) Object {
	var errMsg string
	if left.Type() == ARRAY_OBJ && index.Type() == INTEGER_OBJ {
		array := left.(*Array)
		idx := index.(*Integer)

		length := int64(len(array.Members) - 1)
		if idx.Value > length || idx.Value < -length {
			errMsg = fmt.Sprintf("index %d out of range", idx.Value)
			return &Error{Message: errMsg}
		}
		if idx.Value < 0 {
			return array.Members[(length+1)+idx.Value]
		}
		return array.Members[idx.Value]
	}
	errMsg = fmt.Sprintf("index operator not supported: %s", left.Type())
	return &Error{Message: errMsg}
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

func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalBlockStatements(block []ast.Statement, scope *Scope) Object {
	var results Object

	for _, statement := range block {
		results = Eval(statement, scope)
		if results != nil && results.Type() == RETURN_VALUE_OBJ {
			return results
		}
	}
	return results
}

func evalProgram(program *ast.Program, scope *Scope) Object {
	var results Object

	for _, statement := range program.Statements {
		results = Eval(statement, scope)
		switch s := results.(type) {
		case *ReturnValue:
			return s.Value
		case *Error:
			return s
		}
	}

	return results
}

func evalInfixExpression(operator string, left Object, right Object) Object {
	var errMsg string
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntInfixExpression(operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		errMsg = fmt.Sprintf("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		errMsg = fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	if errMsg != "" {
		return &Error{Message: errMsg}
	}
	return NULL
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

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		if i, ok := right.(*Integer); ok {
			i.Value = -i.Value
			return right
		}
		msg := fmt.Sprintf("unknown operator: %s%s", operator, right.Type())
		return &Error{Message: msg}
	default:
		return NULL
	}
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

func evalArrayLiteral(a *ast.ArrayLiteral, scope *Scope) Object {
	return &Array{Members: evalArgs(a.Members, scope)}
}

func evalHashLiteral(hl *ast.HashLiteral, scope *Scope) Object {
	hashMap := make(map[HashKey]HashPair)
	// TODO: { 1 -> true, 2 -> "five", "three"-> "four } causes some sort of infinite loop
	for key, value := range hl.Pairs {
		key := Eval(key, scope)
		hashable, ok := key.(Hashable)
		if !ok {
			return &Error{Message: fmt.Sprintf("%T not hashable", key.Type())}
		}
		hashMap[hashable.HashKey()] = HashPair{Key: key, Value: Eval(value, scope)}

	}
	return &Hash{Pairs: hashMap}
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

func evalFunctionCall(call *ast.CallExpression, scope *Scope) Object {
	f, ok := scope.Get(call.Function.String())
	if !ok {
		return &Error{Message: fmt.Sprintf("unknown identifier: %s", call.Function.String())}
	}
	fn := f.(*Function)
	fn.Scope = scope
	args := evalArgs(call.Arguments, scope)
	// TODO: If the wrong number of params is passed a panic occurs
	for i, v := range fn.Literal.Parameters {
		scope.Set(v.String(), args[i])
	}
	r := Eval(fn.Literal.Body, scope)
	if obj, ok := r.(*ReturnValue); ok {
		return obj.Value
	}
	return r
}

func evalArgs(args []ast.Expression, scope *Scope) []Object {
	e := []Object{}
	for _, v := range args {
		e = append(e, Eval(v, scope))
	}
	return e
}
