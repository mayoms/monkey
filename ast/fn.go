package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type ArrayLiteral struct {
	Token   token.Token
	Members []Expression
}

func (a *ArrayLiteral) expressionNode()      {}
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteral) String() string {
	var out bytes.Buffer

	members := []string{}
	for _, m := range a.Members {
		members = append(members, m.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(members, ", "))
	out.WriteString("]")
	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(fl.TokenLiteral())
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(" (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString("{ ")
	out.WriteString(fl.Body.String())
	out.WriteString(" }")
	return out.String()
}
