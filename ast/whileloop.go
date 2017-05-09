package ast

import (
	"bytes"
	"monkey/token"
)

type WhileLoop struct {
	Token     token.Token
	Block     *BlockStatement
	Condition Expression
}

func (wl *WhileLoop) expressionNode()      {}
func (wl *WhileLoop) TokenLiteral() string { return wl.Token.Literal }

func (wl *WhileLoop) String() string {
	var out bytes.Buffer

	out.WriteString("while")
	out.WriteString(wl.Condition.String())
	out.WriteString(" { ")
	out.WriteString(wl.Block.String())
	out.WriteString(" }")
	return out.String()
}
