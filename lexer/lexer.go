package lexer

import (
	"unsafe"
)

type Lexer struct {
	input []byte
	index int
	line  int
}

func NewLexer(input []byte) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) NextToken() Token {
	l.skipSpace()
	if l.index >= len(l.input) {
		return Token{Type: EOF, Line: l.line, Literal: ""}
	}
	buff := unsafe.Pointer(&(l.input[l.index]))

	// Two character symbols check
	switch *(*int16)(buff) {
	case 0x3D3D: // == (Little Endian)
		l.index += 2
		return Token{Type: EQ, Line: l.line, Literal: "=="}
	case 0x3D21: // != (Little Endian)
		l.index += 2
		return Token{Type: NOTEQ, Line: l.line, Literal: "!="}
	case 0x3D3C: // <= (Little Endian)
		l.index += 2
		return Token{Type: LTEQ, Line: l.line, Literal: "<="}
	case 0x3D3E: // >= (Little Endian)
		l.index += 2
		return Token{Type: GTEQ, Line: l.line, Literal: ">="}
	}

	// One character symbols check
	switch *(*int8)(buff) {
	case '=':
		l.index += 1
		return Token{Type: ASSIGN, Line: l.line, Literal: "="}
	case '+':
		l.index += 1
		return Token{Type: PLUS, Line: l.line, Literal: "+"}
	case '-':
		l.index += 1
		return Token{Type: MINUS, Line: l.line, Literal: "-"}
	case '!':
		l.index += 1
		return Token{Type: BANG, Line: l.line, Literal: "!"}
	case '/':
		l.index += 1
		return Token{Type: SLASH, Line: l.line, Literal: "/"}
	case '*':
		l.index += 1
		return Token{Type: ASTERISK, Line: l.line, Literal: "*"}
	case '<':
		l.index += 1
		return Token{Type: LT, Line: l.line, Literal: "<"}
	case '>':
		l.index += 1
		return Token{Type: GT, Line: l.line, Literal: ">"}
	case ',':
		l.index += 1
		return Token{Type: COMMA, Line: l.line, Literal: ","}
	case ';':
		l.index += 1
		return Token{Type: SEMICOLON, Line: l.line, Literal: ";"}
	case ':':
		l.index += 1
		return Token{Type: COLON, Line: l.line, Literal: ":"}
	case '(':
		l.index += 1
		return Token{Type: LPAREN, Line: l.line, Literal: "("}
	case ')':
		l.index += 1
		return Token{Type: RPAREN, Line: l.line, Literal: ")"}
	case '{':
		l.index += 1
		return Token{Type: LBRACE, Line: l.line, Literal: "{"}
	case '}':
		l.index += 1
		return Token{Type: RBRACE, Line: l.line, Literal: "}"}
	case '[':
		l.index += 1
		return Token{Type: LBRACKET, Line: l.line, Literal: "["}
	case ']':
		l.index += 1
		return Token{Type: RBRACKET, Line: l.line, Literal: "]"}
	case '"':
		return Token{Type: STRING, Line: l.line, Literal: l.readString()}
	}

	// Multi-character data check (ident, numbers)
	char := *(*byte)(buff)
	switch {
	case isLetter(char):
		ident := l.readIdentifier()
		return Token{
			Type:    LookupIdent(ident),
			Line:    l.line,
			Literal: ident,
		}
	case isDigit(char):
		return Token{
			Type:    INTEGER,
			Line:    l.line,
			Literal: l.readNumber(),
		}
	}

	// In case of unknown character
	l.index++
	return Token{
		Type: ILLEGAL, Line: l.line,
		Literal: string(char),
	}
}

func (l *Lexer) PeekToken() Token {
	// Workaround for cases where we don't
	// want to move the cursor.
	index := l.index
	token := l.NextToken()
	l.index = index
	return token
}

func (l *Lexer) charAt(index int) byte {
	if index < len(l.input) {
		return l.input[index]
	}
	return 0
}

func (l *Lexer) readIdentifier() string {
	start := l.index
	end := l.index
	for isLetter(l.charAt(end)) {
		end++
	}
	l.index = end
	return string(l.input[start:end])
}

func (l *Lexer) readNumber() string {
	start := l.index
	end := l.index
	for isDigit(l.charAt(end)) {
		end++
	}
	l.index = end
	return string(l.input[start:end])
}

func (l *Lexer) readString() string {
	l.index++ // Skip first quotes
	start := l.index
	end := l.index
	for l.charAt(end) != '"' {
		end++
	}
	l.index = end + 1 // Skip second quotes
	return string(l.input[start:end])
}
