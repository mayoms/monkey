package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type StructLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (s *StructLiteral) expressionNode()      {}
func (s *StructLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StructLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range s.Pairs {
		pairs = append(pairs, key.String()+"->"+value.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(")")

	return out.String()

}
