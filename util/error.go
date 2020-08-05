package util

import (
	"bytes"
	"fmt"
	"github.com/Onelio/Eldrlang/lexer"
)

const (
	ExpectedIdent = "expected name declaration but got \"%s\""
	ExpectedParen = "expected opening parenthesis but got \"%s\""
	ExpectedBrace = "expected opening brace but got \"%s\""
	ExpectedCondV = "expected conditional or boolean"
	ExpectedFuncP = "expected %d function parameters"
	UnexpectedEOF = "unexpected end of file, expected \"%s\""
	UnexpectedBRC = "unexpected right brace, expected \"%s\""
	InvalidNumber = "\"%s\" is not a valid number"
	InvalidOpForO = "invalid operator for object"
	InvalidOpComb = "invalid operator combination of objects"
	IllegalLetter = "illegal character \"%s\""
	IllegalOpeAtt = "illegal operation attempt"
	IllegalExprBr = "illegal expresion declaration after break, expected \";\""
	IdentNotFound = "identifier \"%s\" not found"
	IdentNotAFunc = "identifier \"%s\" is not a function"
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
		e.token.Line+1, e.str)
}

type Errors []*Error

func (es *Errors) Len() int {
	return len(*es)
}

func (es *Errors) Add(e *Error) {
	*es = append(*es, e)
}

func (es *Errors) Clear() {
	*es = Errors{}
}

func (es *Errors) String() string {
	var out bytes.Buffer
	for _, e := range *es {
		_, _ = fmt.Println(&out, e.String())
	}
	return out.String()
}
