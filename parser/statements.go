package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return stmt
	}
	p.nextToken()
	stmt.ReturnValue = p.parseExpressionStatement().Expression

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if p.expectPeek(token.IDENT) {
		stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if p.expectPeek(token.ASSIGN) {
		p.nextToken()
		stmt.Value = p.parseExpressionStatement().Expression
	}

	return stmt
}

func (p *Parser) parseIncludeStatement() *ast.IncludeStatement {
	stmt := &ast.IncludeStatement{Token: p.curToken}

	if p.expectPeek(token.IDENT) {
		stmt.ImportFile = p.parseExpressionStatement().Expression
	}
	return stmt
}
