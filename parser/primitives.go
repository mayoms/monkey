package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"strconv"
)

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseIdentifier() ast.Expression {
	identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.peekTokenIs(token.LPAREN) {
		return identifier
	}
	callexp := p.parseCallExpressions(identifier)
	return callexp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}
