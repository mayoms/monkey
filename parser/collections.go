package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseIndexExpression(arr ast.Expression) ast.Expression {
	indexExp := &ast.IndexExpression{Token: p.curToken, Left: arr}
	p.nextToken()
	index := p.parseExpression(LOWEST)
	if p.expectPeek(token.RBRACKET) {
		indexExp.Index = index
	}
	return indexExp
}

func (p *Parser) parseHashExpression() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return hash
	}
	for !p.curTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if p.expectPeek(token.ARROW) {
			p.nextToken()
			hash.Pairs[key] = p.parseExpression(LOWEST)
			p.nextToken()
		}
	}
	return hash
}

func (p *Parser) parseArrayExpression() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken, Members: []ast.Expression{}}
	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
		return array
	}
	for !p.curTokenIs(token.RBRACKET) {
		p.nextToken()
		array.Members = append(array.Members, p.parseExpression(LOWEST))
		p.nextToken()
	}
	return array
}
