package token

type TokenType int

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = iota

	EOF

	// Identifiers + literals
	IDENT  // add, foobar, x, y, ...
	NUMBER // 1343456

	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	STAR
	SLASH
	LT
	GT

	EQ
	NOT_EQ

	// Delimiters
	COMMA
	SEMI
	LPAREN
	RPAREN
	LBRACE
	RBRACE

	// Keywords
	FN
	LET
	TRUE
	FALSE
	IF
	ELSE
	RETURN
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF: "EOF",

	// Identifiers + literals
	IDENT:  "IDENT",  // add, foobar, x, y, ...
	NUMBER: "NUMBER", // 1343456

	// Operators
	ASSIGN: "ASSIGN",
	PLUS:   "PLUS",
	MINUS:  "MINUS",
	BANG:   "BANG",
	STAR:   "STAR",
	SLASH:  "SLASH",
	LT:     "LT",
	GT:     "GT",

	EQ:     "EQ",
	NOT_EQ: "NOT_EQ",

	// Delimiters
	COMMA:  "COMMA",
	SEMI:   "SEMI",
	LPAREN: "LPAREN",
	RPAREN: "RPAREN",
	LBRACE: "LBRACE",
	RBRACE: "RBRACE",

	// Keywords
	FN:     "FN",
	LET:    "LET",
	TRUE:   "TRUE",
	FALSE:  "FALSE",
	IF:     "IF",
	ELSE:   "ELSE",
	RETURN: "RETURN",
}

var keywords = map[string]TokenType{
	"fn":     FN,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tt, ok := keywords[ident]; ok {
		return tt
	}
	return IDENT
}

func (tt TokenType) String() string {
	return tokens[tt]
}
