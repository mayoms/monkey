package parser

import (
	"bytes"
	"fmt"
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
		fmt.Println(p.curToken)
		if p.peekTokenIs(token.EOF) {
			p.noPrefixParseFnError(p.peekToken.Type)
			break
		}
		if p.curTokenIs(token.LBRACE) {
			out.WriteString("{")
			p.nextToken()
			if p.curTokenIs(token.RBRACE) {
				out.WriteString("}")
				p.nextToken()
				continue
			}
			if !p.curTokenIs(token.BYTES) {
				exp := p.parseExpression(LOWEST)
				out.WriteString(exp.String())
				is.ExprList = append(is.ExprList, exp)
				p.nextToken()
			}
		}
		out.WriteString(p.curToken.Literal)
		p.nextToken()
	}
	fmt.Println(out.String())
	is.Value = out.String()
	return is
}
