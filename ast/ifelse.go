package ast

import (
	"bytes"
	"monkey/token"
)

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifex *IfExpression) expressionNode()      {}
func (ifex *IfExpression) TokenLiteral() string { return ifex.Token.Literal }

func (ifex *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString("(")
	out.WriteString(ifex.Condition.String())
	out.WriteString(")")
	out.WriteString(" { ")
	out.WriteString(ifex.Consequence.String())
	out.WriteString(" } ")
	if ifex.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(" { ")
		out.WriteString(ifex.Alternative.String())
		out.WriteString(" } ")
	}

	return out.String()
}
