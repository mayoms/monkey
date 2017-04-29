package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	fn.Parameters = p.parseExpressionArray(fn.Parameters, token.RPAREN)
	if p.expectPeek(token.LBRACE) {
		fn.Body = p.parseBlockStatement().(*ast.BlockStatement)
	}
	return fn
}

func (p *Parser) parseCallExpressions(f ast.Expression) ast.Expression {
	call := &ast.CallExpression{Token: p.curToken, Function: f}
	call.Arguments = p.parseExpressionArray(call.Arguments, token.RPAREN)
	return call
}

func (p *Parser) parseExpressionArray(a []ast.Expression, closure token.TokenType) []ast.Expression {
	if p.peekTokenIs(closure) {
		p.nextToken()
		return a
	}
	p.nextToken()
	a = append(a, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		a = append(a, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(closure) {
		return nil
	}
	return a
}

func (p *Parser) parseMethodCallExpression(obj ast.Expression) ast.Expression {
	methodCall := &ast.MethodCallExpression{Token: p.curToken, Object: obj}
	p.nextToken()
	name := p.parseIdentifier()
	if !p.peekTokenIs(token.LPAREN) {
		methodCall.Call = p.parseExpression(LOWEST)
	} else {
		p.nextToken()
		methodCall.Call = p.parseCallExpressions(name)
	}
	return methodCall
}
