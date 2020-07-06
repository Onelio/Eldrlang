package parser

import (
	"Eldrlang/lexer"
	"bytes"
	"strings"
)

type Assign struct {
	Token lexer.Token
	Left  Expression
	Right Expression
}

func (a *Assign) Literal() string { return a.Token.Literal }
func (a *Assign) String() string {
	var out bytes.Buffer
	out.WriteString(a.Left.String())
	out.WriteString(" = ")
	out.WriteString(a.Right.String())
	return out.String()
}

func (p *Parser) newAssign(left Expression) Expression {
	expression := &Assign{
		Token: p.token,
		Left:  left,
	}
	var exp Expression
	for !p.isPeekToken(lexer.SEMICOLON) {
		p.nextToken() // Go next first to skip = opcode
		exp = p.parseToken(exp)
	}
	expression.Right = exp
	return expression
}

type Prefix struct {
	Token    lexer.Token
	Operator string
	Right    Expression
}

func (p *Prefix) Literal() string { return p.Token.Literal }
func (p *Prefix) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

func (p *Parser) newPrefix() Expression {
	prefix := &Prefix{
		Token:    p.token,
		Operator: p.token.Literal,
	}
	p.nextToken()
	prefix.Right = p.parseToken(nil)
	return prefix
}

type Infix struct {
	Token    lexer.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *Infix) Literal() string { return i.Token.Literal }
func (i *Infix) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

func (p *Parser) newInfix(left Expression) Expression {
	expression := &Infix{
		Token:    p.token,
		Operator: p.token.Literal,
		Left:     left,
	}
	p.nextToken()
	expression.Right = p.parseToken(nil)
	return expression
}

type FuncCall struct {
	Token     lexer.Token
	Function  Expression
	Arguments []Expression
}

func (fc *FuncCall) Literal() string { return fc.Token.Literal }
func (fc *FuncCall) String() string {
	var out bytes.Buffer
	var args []string
	for _, a := range fc.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(fc.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

func (p *Parser) newFuncCall(function Expression) Expression {
	exp := &FuncCall{Token: p.token, Function: function}
	exp.Arguments = p.newParameters()
	return exp
}

type Return struct {
	Token lexer.Token
	Exp   Expression
}

func (r *Return) Literal() string { return r.Token.Literal }
func (r *Return) String() string {
	var out bytes.Buffer
	out.WriteString("return")
	out.WriteString(r.Exp.String())
	return out.String()
}

func (p *Parser) newReturn() Expression {
	expression := &Return{
		Token: p.token,
	}
	p.nextToken()
	expression.Exp = p.parseExpression()
	return expression
}
