package lexer

const (
	EOF = 1 << iota
	IDENT

	// Types
	BOOLEAN
	INTEGER
	STRING

	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	SLASH
	ASTERISK

	LT
	LTEQ
	GT
	GTEQ
	EQ
	NOTEQ

	// Delimiters
	COMMA
	SEMICOLON
	COLON

	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET

	// Reserved Keywords
	VARIABLE
	FUNCTION
	RETURN
	IMPORT
	IF
	ELSE
	FOR
)

type Type int

var keywords = map[string]Type{
	"var":    VARIABLE,
	"fun":    FUNCTION,
	"return": RETURN,
	"true":   BOOLEAN,
	"false":  BOOLEAN,
	"import": IMPORT,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
}

type Token struct {
	Type
	Line    int
	Literal string
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipSpace() {
	for l.index < len(l.Input) {
		char := l.Input[l.index]
		switch {
		case char == '\n':
			l.line++
			fallthrough
		case char <= ' ':
			l.index++
			continue
		}
		return
	}
}