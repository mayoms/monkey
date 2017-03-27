package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseFunctionLiteral() ast.Expression {

	lit := &ast.FunctionLiteral{Token: p.curToken}

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
		} else {
			return nil
		}
	}
	lit.Parameters = params

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement().(*ast.BlockStatement)

	return lit
}
