package eval

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue)
		if value != nil {
			return &object.ReturnValue{Value: value}
		}
		return NULL
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		if left.Type() == object.ERROR_OBJ {
			return left
		} else if right.Type() == object.ERROR_OBJ {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if right.Type() == object.ERROR_OBJ {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.IfExpression:
		condition := Eval(node.Condition)
		if isTrue(condition) {
			return evalBlockStatements(node.Consequence.Statements)
		} else if node.Alternative != nil {
			return evalBlockStatements(node.Alternative.Statements)
		}
		return NULL
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}
	return nil
}

func isTrue(obj object.Object) bool {
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

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalBlockStatements(block []ast.Statement) object.Object {
	var results object.Object

	for _, statement := range block {
		results = Eval(statement)
		if results != nil && results.Type() == object.RETURN_VALUE_OBJ {
			return results
		}
	}
	return results
}

func evalProgram(program *ast.Program) object.Object {
	var results object.Object

	for _, statement := range program.Statements {
		results = Eval(statement)
		switch s := results.(type) {
		case *object.ReturnValue:
			return s.Value
		case *object.Error:
			return s
		}
	}

	return results
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	var errMsg string
	switch {
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		errMsg = fmt.Sprintf("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case operator == "+":
		errMsg = fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	if errMsg != "" {
		return &object.Error{Message: errMsg}
	}
	return NULL
}

func evalIntInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	l := left.(*object.Integer)
	r := right.(*object.Integer)

	switch operator {
	case "+":
		return &object.Integer{Value: l.Value + r.Value}
	case "-":
		return &object.Integer{Value: l.Value - r.Value}
	case "*":
		return &object.Integer{Value: l.Value * r.Value}
	case "/":
		return &object.Integer{Value: l.Value / r.Value}
	case ">":
		return nativeBoolToBooleanObject(l.Value > r.Value)
	case "<":
		return nativeBoolToBooleanObject(l.Value < r.Value)
	}
	return NULL
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		if i, ok := right.(*object.Integer); ok {
			i.Value = -i.Value
			return right
		}
		msg := fmt.Sprintf("unknown operator: %s%s", operator, right.Type())
		return &object.Error{Message: msg}
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
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
