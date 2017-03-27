package parser

import (
	"go/token"
	"monkey/ast"
)

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement().(*ast.BlockStatement)
	if p.expectPeek(token.ELSE) && p.expectPeek(token.LBRACE) {
		expression.Alternative = p.parseBlockStatement().(*ast.BlockStatement)
	}
	return expression
}
