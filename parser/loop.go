package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseDoLoopExpression() ast.Expression {
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	loop := &ast.DoLoop{Token: p.curToken}
	p.expectPeek(token.LBRACE)
	loop.Block = p.parseBlockStatement().(*ast.BlockStatement)
	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	return loop
}
