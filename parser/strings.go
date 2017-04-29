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
		out.WriteString(p.curToken.Literal)
		p.nextToken()
	}
	is.Value = out.String()
	return is
}
