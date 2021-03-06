package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	EQ       = "=="
	NEQ      = "!="
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MOD      = "%"

	LT        = "<"
	GT        = ">"
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
	COLON    = ":"
	DOT      = "."
	ARROW    = "->"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	INCLUDE  = "INCLUDE"
	STRING   = "STRING"
	ISTRING  = "ISTRING"
	BYTES    = "BYTES"
	AND      = "AND"
	OR       = "OR"
	STRUCT   = "STRUCT"
	DO       = "DO"
	BREAK    = "BREAK"
)

var keywords = map[string]TokenType{
	"fn":      FUNCTION,
	"let":     LET,
	"true":    TRUE,
	"false":   FALSE,
	"if":      IF,
	"else":    ELSE,
	"return":  RETURN,
	"include": INCLUDE,
	"and":     AND,
	"or":      OR,
	"struct":  STRUCT,
	"do":      DO,
	"break":   BREAK,
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
