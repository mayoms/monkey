package parser

import (
	"monkey/ast"
	"monkey/token"
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
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if p.expectPeek(token.LBRACE) {
			expression.Alternative = p.parseBlockStatement().(*ast.BlockStatement)
		}
	}

	return expression
}

func (p *Parser) parseDoLoopExpression() ast.Expression {
	loop := &ast.DoLoop{Token: p.curToken}
	p.expectPeek(token.LBRACE)
	loop.Block = p.parseBlockStatement().(*ast.BlockStatement)
	return loop
}
