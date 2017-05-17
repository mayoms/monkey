package lexer

import (
	"bytes"
	"errors"
	"monkey/token"
)

type Lexer struct {
	input        string
	ch           byte
	position     int
	readPosition int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

var tokenMap = map[byte]token.TokenType{
	'=': token.EQ,
	'.': token.DOT,
	';': token.SEMICOLON,
	'(': token.LPAREN,
	')': token.RPAREN,
	'{': token.LBRACE,
	'}': token.RBRACE,
	'[': token.LBRACKET,
	']': token.RBRACKET,
	'+': token.PLUS,
	',': token.COMMA,
	'-': token.MINUS,
	'!': token.BANG,
	'*': token.ASTERISK,
	'/': token.SLASH,
	'<': token.LT,
	'>': token.GT,
	':': token.COLON,
	'%': token.MOD,
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	if t, ok := tokenMap[l.ch]; ok {
		switch t {
		case token.EQ:
			if l.peekChar() == '=' {
				tok = token.Token{token.EQ, string(l.ch) + string(l.ch)}
				l.readChar()
			} else {
				tok = newToken(token.ASSIGN, l.ch)
			}
		case token.MINUS:
			if l.peekChar() == '>' {
				tok = token.Token{token.ARROW, string(l.ch) + string(l.peekChar())}
				l.readChar()
			} else {
				tok = newToken(token.MINUS, l.ch)
			}
		case token.BANG:
			if l.peekChar() == '=' {
				tok = token.Token{token.NEQ, string(l.ch) + string(l.peekChar())}
				l.readChar()
			} else {
				tok = newToken(token.BANG, l.ch)
			}
		default:
			tok = newToken(t, l.ch)
		}
		l.readChar()
		return tok
	}
	return l.readBytesToken()
}

func (l *Lexer) readBytesToken() token.Token {
	var tok token.Token
	switch {
	case l.ch == 0:
		tok.Literal = ""
		tok.Type = token.EOF
		return tok
	case isLetter(l.ch):
		tok.Literal = l.readIdentifier()
		tok.Type = token.LookupIdent(tok.Literal)
		return tok
	case isDigit(l.ch):
		tok.Type = token.INT
		tok.Literal = l.readNumber()
		return tok
	case isQuote(l.ch):
		if s, err := l.readString(); err == nil {
			tok.Type = token.STRING
			tok.Literal = s
			return tok
		}
	case isSingleQuote(l.ch):
		if s, err := l.readInterpString(); err == nil {
			tok.Type = token.ISTRING
			tok.Literal = s
			return tok
		}
	}
	l.readChar()
	return newToken(token.ILLEGAL, l.ch)
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readString() (string, error) {
	start := l.position + 1
	for {
		l.readChar()
		if isQuote(l.ch) {
			l.readChar()
			break
		}
		if l.ch == 0 {
			err := errors.New("")
			return "", err
		}
	}
	return l.input[start : l.position-1], nil
}

func (l *Lexer) readInterpString() (string, error) {
	start := l.position + 1
	var out bytes.Buffer
	pos := "0"[0]
	for {
		l.readChar()
		if isSingleQuote(l.ch) {
			l.readChar()
			break
		}
		if l.ch == 0 {
			err := errors.New("")
			return "", err
		}
		if l.ch == 123 {
			if l.peekChar() != 125 {
				out.WriteByte(l.ch)
				for l.ch != 125 || l.ch == 0 {
					l.readChar()
				}
				if l.ch != 0 {
					out.WriteByte(pos)
					pos++
				}
			}
		}
		out.WriteByte(l.ch)
	}
	l.position = start - 1
	l.readPosition = start
	l.ch = l.input[start]
	return out.String(), nil
}

func (l *Lexer) NextInterpToken() token.Token {
	var tok token.Token
	for {
		if l.ch == '{' {
			if l.peekChar() == '}' {
				continue
			}
			tok = newToken(token.LBRACE, l.ch)
			break
		}
		if l.ch == 0 {
			tok.Type = token.EOF
			tok.Literal = ""
			break
		}
		if isSingleQuote(l.ch) {
			tok = newToken(token.ISTRING, l.ch)
			break
		}
		l.readChar()
	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isQuote(ch byte) bool {
	return ch == 34
}

func isSingleQuote(ch byte) bool {
	return ch == 39
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
