package ast

import (
	"bytes"
	"monkey/token"
)

type DoLoop struct {
	Token token.Token
	Block *BlockStatement
}

func (dl *DoLoop) expressionNode()      {}
func (dl *DoLoop) TokenLiteral() string { return dl.Token.Literal }

func (dl *DoLoop) String() string {
	var out bytes.Buffer

	out.WriteString("do")
	out.WriteString(" { ")
	out.WriteString(dl.Block.String())
	out.WriteString(" }")
	return out.String()
}
