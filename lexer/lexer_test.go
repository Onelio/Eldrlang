package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := []byte(`var five = 5;
var ten = 10;

var add = fun(x, y) {
  x + y;
};

var result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}`)

	tests := []struct {
		expectedType    Type
		expectedLiteral string
	}{
		{VARIABLE, "var"},
		{IDENT, "five"},
		{ASSIGN, "="},
		{INTEGER, "5"},
		{SEMICOLON, ";"},
		{VARIABLE, "var"},
		{IDENT, "ten"},
		{ASSIGN, "="},
		{INTEGER, "10"},
		{SEMICOLON, ";"},
		{VARIABLE, "var"},
		{IDENT, "add"},
		{ASSIGN, "="},
		{FUNCTION, "fun"},
		{LPAREN, "("},
		{IDENT, "x"},
		{COMMA, ","},
		{IDENT, "y"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{IDENT, "x"},
		{PLUS, "+"},
		{IDENT, "y"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{SEMICOLON, ";"},
		{VARIABLE, "var"},
		{IDENT, "result"},
		{ASSIGN, "="},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "five"},
		{COMMA, ","},
		{IDENT, "ten"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{BANG, "!"},
		{MINUS, "-"},
		{SLASH, "/"},
		{ASTERISK, "*"},
		{INTEGER, "5"},
		{SEMICOLON, ";"},
		{INTEGER, "5"},
		{LT, "<"},
		{INTEGER, "10"},
		{GT, ">"},
		{INTEGER, "5"},
		{SEMICOLON, ";"},
		{IF, "if"},
		{LPAREN, "("},
		{INTEGER, "5"},
		{LT, "<"},
		{INTEGER, "10"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{BOOLEAN, "true"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{ELSE, "else"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{BOOLEAN, "false"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{INTEGER, "10"},
		{EQ, "=="},
		{INTEGER, "10"},
		{SEMICOLON, ";"},
		{INTEGER, "10"},
		{NOTEQ, "!="},
		{INTEGER, "9"},
		{SEMICOLON, ";"},
		{STRING, "foobar"},
		{STRING, "foo bar"},
		{LBRACKET, "["},
		{INTEGER, "1"},
		{COMMA, ","},
		{INTEGER, "2"},
		{RBRACKET, "]"},
		{SEMICOLON, ";"},
		{LBRACE, "{"},
		{STRING, "foo"},
		{COLON, ":"},
		{STRING, "bar"},
		{RBRACE, "}"},
		{EOF, ""},
	}

	l := Lexer{Input: input}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%d, got=%d",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
