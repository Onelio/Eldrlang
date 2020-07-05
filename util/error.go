package util

import (
	"Eldrlang/lexer"
	"bytes"
	"fmt"
)

const (
	ExpectedIdent = "expected name declaration but got \"%s\""
	ExpectedParen = "expected opening parenthesis but got \"%s\""
	ExpectedBrace = "expected opening brace but got \"%s\""
	ExpectedInitL = "expected initialization for loop variable"
	ExpectedCondL = "expected condition for loop action"
	ExpectedLoopC = "expected loop closing parenthesis"
	UnexpectedEOF = "unexpected end of file, expected \"%s\""
	InvalidNumber = "\"%s\" is not a valid number"
	IllegalLetter = "illegal character \"%s\""
	IllegalOpeAtt = "illegal operation attempt"
)

type Error struct {
	token lexer.Token
	str   string
}

func NewError(t lexer.Token, f string, a ...interface{}) *Error {
	return &Error{
		token: t,
		str:   fmt.Sprintf(f, a...),
	}
}

func (e *Error) String() string {
	return fmt.Sprintf("* Error at L%d %s",
		e.token.Line, e.str)
}

type Errors []*Error

func (es *Errors) Len() int {
	return len(*es)
}

func (es *Errors) Add(e *Error) {
	*es = append(*es, e)
}

func (es *Errors) String() string {
	var out bytes.Buffer
	for _, e := range *es {
		_, _ = fmt.Fprintln(&out, e.String())
	}
	return out.String()
}
