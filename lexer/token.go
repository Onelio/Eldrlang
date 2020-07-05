package lexer

const (
	EOF = iota
	IDENT
	ILLEGAL

	// Types
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
	TRUE
	FALSE
	IMPORT
	IF
	ELSE
	LOOP
)

type Type int

var keywords = map[string]Type{
	"var":    VARIABLE,
	"fun":    FUNCTION,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"import": IMPORT,
	"if":     IF,
	"else":   ELSE,
	"loop":   LOOP,
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
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipSpace() {
	for l.index < len(l.input) {
		char := l.input[l.index]
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
