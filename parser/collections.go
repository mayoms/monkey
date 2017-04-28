package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseSliceExpression(start ast.Expression) ast.Expression {
	slice := &ast.SliceExpression{Token: p.curToken, StartIndex: start}
	if p.peekTokenIs(token.RBRACKET) {
		slice.EndIndex = nil
	} else {
		p.nextToken()
		slice.EndIndex = p.parseExpression(LOWEST)
	}
	return slice
}

func (p *Parser) parseIndexExpression(arr ast.Expression) ast.Expression {
	var index ast.Expression
	indexExp := &ast.IndexExpression{Token: p.curToken, Left: arr}
	if p.peekTokenIs(token.COLON) {
		indexTok := token.Token{Type: token.INT, Literal: "0"}
		prefix := &ast.IntegerLiteral{Token: indexTok, Value: int64(0)}
		p.nextToken()
		index = p.parseSliceExpression(prefix)
	} else {
		p.nextToken()
		index = p.parseExpression(LOWEST)
	}
	indexExp.Index = index
	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
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
		if !p.expectPeek(token.ARROW) {
			return nil
		}
		p.nextToken()
		hash.Pairs[key] = p.parseExpression(LOWEST)
		p.nextToken()
	}
	return hash
}

func (p *Parser) parseStructExpression() ast.Expression {
	s := &ast.StructLiteral{Token: p.curToken}
	s.Pairs = make(map[ast.Expression]ast.Expression)
	p.expectPeek(token.LPAREN)
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return s
	}
	for !p.curTokenIs(token.RPAREN) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.ARROW) {
			return nil
		}
		p.nextToken()
		s.Pairs[key] = p.parseExpression(LOWEST)
		p.nextToken()
	}
	return s
}

func (p *Parser) parseArrayExpression() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Members = p.parseExpressionArray(array.Members, token.RBRACKET)
	return array
}
