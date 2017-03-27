package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseFunctionLiteral() ast.Expression {

	expression := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	params := []*ast.Identifier{}

	for p.curToken.Type != token.RPAREN {
		if p.curTokenIs(token.IDENT) {
			ident := p.parseIdentifier().(*ast.Identifier)
			params = append(params, ident)
			p.nextToken()
		} else if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	expression.Parameters = params

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement().(*ast.BlockStatement)

	return expression
}
