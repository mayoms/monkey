package parser

import (
	"bytes"
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseStringLiteralExpression() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseInterpolatedString() ast.Expression {
	is := &ast.InterpolatedString{Token: p.curToken, ExprList: []ast.Expression{}}

	p.nextToken()
	var out bytes.Buffer

	for !p.curTokenIs(token.ISTRING) {
		if !p.curTokenIs(token.BYTES) {
			out.WriteString("{")
			exp := p.parseExpression(LOWEST)
			out.WriteString(exp.String())
			is.ExprList = append(is.ExprList, exp)
			p.expectPeek(token.RBRACE)
		}
		out.WriteString(p.curToken.Literal)
		p.nextToken()
	}
	is.Value = out.String()
	return is
}
