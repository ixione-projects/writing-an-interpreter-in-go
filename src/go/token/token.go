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
	ASSIGN: "=",
	PLUS:   "+",
	MINUS:  "-",
	BANG:   "!",
	STAR:   "*",
	SLASH:  "/",
	LT:     "<",
	GT:     ">",

	EQ:     "==",
	NOT_EQ: "!=",

	// Delimiters
	COMMA:  ",",
	SEMI:   ";",
	LPAREN: "(",
	RPAREN: ")",
	LBRACE: "{",
	RBRACE: "}",

	// Keywords
	FN:     "fn",
	LET:    "let",
	TRUE:   "true",
	FALSE:  "false",
	IF:     "if",
	ELSE:   "else",
	RETURN: "return",
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
