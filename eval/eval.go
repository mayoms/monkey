package eval

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	default:
		fmt.Printf("%T", node)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var results object.Object

	for _, statement := range stmts {
		results = Eval(statement)
	}

	return results
}
