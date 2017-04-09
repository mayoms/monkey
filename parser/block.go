package parser

import (
	"monkey/ast"
	"monkey/token"
)

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
