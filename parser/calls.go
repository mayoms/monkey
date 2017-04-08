package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseFunctionLiteral() ast.Expression {

	fn := &ast.FunctionLiteral{Token: p.curToken, Parameters: []*ast.Identifier{}}
	if p.expectPeek(token.LPAREN) {
		p.nextToken()
	}

	for !p.curTokenIs(token.RPAREN) {
		fn.Parameters = append(fn.Parameters, p.parseIdentifier().(*ast.Identifier))
		p.nextToken()
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	if p.expectPeek(token.LBRACE) {
		fn.Body = p.parseBlockStatement().(*ast.BlockStatement)
	}
	return fn
}

func (p *Parser) parseCallExpressions(f ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: p.curToken, Function: f}
	expression.Arguments = []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return expression
	}
	for !p.curTokenIs(token.RPAREN) {
		p.nextToken()
		expression.Arguments = append(expression.Arguments, p.parseExpression(LOWEST))
		p.nextToken()
	}

	return expression
}

func (p *Parser) parseMethodCallExpressions(obj ast.Expression) ast.Expression {
	methodCall := &ast.MethodCallExpression{Token: p.curToken, Object: obj}
	p.nextToken()
	methodCall.Call = p.parseExpression(LOWEST)
	return methodCall
}

func (p *Parser) parseBlockStatement() ast.Expression {
	expression := &ast.BlockStatement{Token: p.curToken}
	expression.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) {
		stmt := p.parseStatement()
		if stmt != nil {
			expression.Statements = append(expression.Statements, stmt)
		}
		if p.peekTokenIs(token.EOF) {
			break
		}
		p.nextToken()
	}
	return expression
}
