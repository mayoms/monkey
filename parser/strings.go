package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseStringLiteralExpression() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseInterpolatedString() ast.Expression {
	is := &ast.InterpolatedString{Token: p.curToken, Value: p.curToken.Literal, ExprMap: make(map[byte]ast.Expression)}

	key := "0"[0]
	for {
		if p.curTokenIs(token.LBRACE) {
			p.nextToken()
			expr := p.parseExpression(LOWEST)
			is.ExprMap[key] = expr
			key++
		}
		p.nextInterpToken()
		if p.curTokenIs(token.ISTRING) {
			break
		}
	}
	return is
}
