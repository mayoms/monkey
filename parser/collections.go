package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseIndexExpression(arr ast.Expression) ast.Expression {
	indexExp := &ast.IndexExpression{Token: p.curToken, Left: arr}
	p.nextToken()
	index := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	indexExp.Index = index
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
		if !p.expectPeek(token.ARROW) {
			return nil
		}
		p.nextToken()
		hash.Pairs[key] = p.parseExpression(LOWEST)
		p.nextToken()
	}
	return hash
}

func (p *Parser) parseArrayExpression() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Members = p.parseExpressionArray(array.Members, token.RBRACKET)
	return array
}
