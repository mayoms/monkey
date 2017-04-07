package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (h *HashLiteral) expressionNode()      {}
func (h *HashLiteral) TokenLiteral() string { return h.Token.Literal }
func (h *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range h.Pairs {
		pairs = append(pairs, key.String()+"->"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()

}
