package parser

import (
	"fmt"
	"io/ioutil"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"os"
)

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return stmt
	}
	p.nextToken()
	stmt.ReturnValue = p.parseExpressionStatement().Expression

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if p.expectPeek(token.IDENT) {
		stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if p.expectPeek(token.ASSIGN) {
		p.nextToken()
		stmt.Value = p.parseExpressionStatement().Expression
	}

	return stmt
}

func (p *Parser) parseIncludeStatement() *ast.IncludeStatement {
	stmt := &ast.IncludeStatement{Token: p.curToken}

	if p.expectPeek(token.IDENT) {
		stmt.IncludePath = p.parseExpressionStatement().Expression
	}
	program, module, err := p.getIncludedStatements(stmt.IncludePath.String())
	if err != nil {
		p.errors = append(p.errors, err.Error())
	}
	stmt.Program = program
	stmt.IsModule = module
	return stmt
}

func (p *Parser) getIncludedStatements(importpath string) (*ast.Program, bool, error) {
	module := false
	path := p.path
	f, err := ioutil.ReadFile(path + "/" + importpath + ".my")
	if err != nil {
		path = path + "/" + importpath
		_, err := os.Stat(path)
		if err != nil {
			return nil, module, fmt.Errorf("no file or directory: %s.my, %s", importpath, path)
		}
		m, err := ioutil.ReadFile(path + "/module.my")
		if err != nil {
			return nil, module, err
		}
		module = true
		f = m
	}
	l := lexer.New(string(f))
	ps := New(l, path)
	parsed := ps.ParseProgram()
	if len(ps.errors) != 0 {
		p.errors = append(p.errors, ps.errors...)
	}
	return parsed, module, nil
}
