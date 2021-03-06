package parser

import (
	"bytes"
	"github.com/Onelio/Eldrlang/lexer"
	"github.com/Onelio/Eldrlang/util"
	"strconv"
)

type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) Literal() string { return i.Token.Literal }
func (i *Identifier) String() string  { return i.Value }

func (p *Parser) newIdentifier() Expression {
	ident := &Identifier{Token: p.token, Value: p.token.Literal}
	return ident
}

func (p *Parser) newParameters() []Expression {
	var identifiers []Expression
	p.nextToken() // Skip opening "("
	var exp Expression
	for !p.isTokenOrEOF(lexer.RPAREN) {
		exp = p.parseToken(exp)
		p.nextToken()
		if p.isToken(lexer.COMMA) {
			identifiers = append(identifiers, exp)
			exp = nil
			p.nextToken()
		}
	}
	if exp != nil { // Cause when exiting there may be still one left
		identifiers = append(identifiers, exp)
	}
	if p.isToken(lexer.EOF) {
		err := util.NewError(p.token, util.UnexpectedEOF, ")")
		p.errors.Add(err)
		return nil
	}
	return identifiers
}

type Boolean struct {
	Token lexer.Token
	Value bool
}

func (b *Boolean) Literal() string { return b.Token.Literal }
func (b *Boolean) String() string  { return b.Token.Literal }

func (p *Parser) newBoolean() Expression {
	return &Boolean{Token: p.token, Value: p.isToken(lexer.TRUE)}
}

type Integer struct {
	Token lexer.Token
	Value int64
}

func (i *Integer) Literal() string { return i.Token.Literal }
func (i *Integer) String() string  { return i.Token.Literal }

func (p *Parser) newInteger() Expression {
	integer := &Integer{Token: p.token}

	value, err := strconv.ParseInt(p.token.Literal, 0, 64)
	if err != nil {
		err := util.NewError(p.token, util.InvalidNumber, p.token.Literal)
		p.errors.Add(err)
		return nil
	}
	integer.Value = value
	return integer
}

type String struct {
	Token lexer.Token
	Value string
}

func (s *String) Literal() string { return s.Token.Literal }
func (s *String) String() string {
	var out bytes.Buffer
	out.WriteString("\"")
	out.WriteString(s.Token.Literal)
	out.WriteString("\"")
	return out.String()
}

func (p *Parser) newString() Expression {
	return &String{Token: p.token, Value: p.token.Literal}
}
