package parser

import (
	"bytes"
	"github.com/Onelio/Eldrlang/lexer"
	"github.com/Onelio/Eldrlang/util"
	"strings"
)

type Variable struct {
	Token lexer.Token
	Name  *Identifier
}

func (d *Variable) Literal() string { return d.Name.String() }
func (d *Variable) String() string {
	var out bytes.Buffer
	out.WriteString(d.Token.Literal + " ")
	out.WriteString(d.Name.String())
	return out.String()
}

func (p *Parser) newVariable() Expression {
	def := &Variable{Token: p.token}

	if p.nextToken() != lexer.IDENT {
		err := util.NewError(p.token, util.ExpectedIdent, p.token.Literal)
		p.errors.Add(err)
		return nil
	}
	def.Name = p.newIdentifier().(*Identifier)

	if p.nextToken() != lexer.SEMICOLON {
		return p.parseToken(def)
	}
	return def
}

type Block struct {
	Token lexer.Token
	Nodes []Node
}

func (b *Block) Literal() string { return b.Token.Literal }
func (b *Block) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range b.Nodes {
		out.WriteString("\t" + s.String() + ";\n")
	}
	out.WriteString("}")
	return out.String()
}

func (p *Parser) newBlock() *Block {
	block := &Block{Token: p.token}
	p.nextToken() // Skip { opening
	for !p.isTokenOrEOF(lexer.RBRACE) {
		node := p.parseStatement()
		if node != nil {
			block.Nodes = append(block.Nodes, node)
		}
		p.nextToken()
	}
	if p.isToken(lexer.EOF) {
		err := util.NewError(p.token, util.UnexpectedEOF, "}")
		p.errors.Add(err)
		return nil
	}
	return block
}

type Conditional struct {
	Token   lexer.Token
	Require Expression
	To      *Block
	Else    *Block
}

func (c *Conditional) Literal() string { return c.Token.Literal }
func (c *Conditional) String() string {
	var out bytes.Buffer
	out.WriteString("if (")
	out.WriteString(c.Require.String())
	out.WriteString(") ")
	out.WriteString(c.To.String())

	if c.Else != nil {
		out.WriteString(" else ")
		out.WriteString(c.Else.String())
	}

	return out.String()
}

func (p *Parser) newConditional() Statement {
	cond := &Conditional{Token: p.token}

	if p.nextToken() != lexer.LPAREN {
		err := util.NewError(p.token, util.ExpectedParen, p.token.Literal)
		p.errors.Add(err)
		return nil
	}
	cond.Require = p.parseGroupExpression()

	if p.nextToken() != lexer.LBRACE {
		err := util.NewError(p.token, util.ExpectedBrace, p.token.Literal)
		p.errors.Add(err)
		return nil
	}
	cond.To = p.newBlock()

	if p.isPeekToken(lexer.ELSE) {
		p.nextToken() // Skip else token
		if p.nextToken() != lexer.LBRACE {
			err := util.NewError(p.token, util.ExpectedBrace, p.token.Literal)
			p.errors.Add(err)
			return nil
		}
		cond.Else = p.newBlock()
	}
	return cond
}

type Loop struct {
	Token lexer.Token
	Body  *Block
}

func (w *Loop) Literal() string { return w.Token.Literal }
func (w *Loop) String() string {
	var out bytes.Buffer
	out.WriteString("loop ")
	out.WriteString(w.Body.String())
	return out.String()
}

func (p *Parser) newLoop() Statement {
	while := &Loop{Token: p.token}
	if p.nextToken() != lexer.LBRACE {
		err := util.NewError(p.token, util.ExpectedBrace, p.token.Literal)
		p.errors.Add(err)
		return nil
	}
	while.Body = p.newBlock()
	return while
}

type Function struct {
	Token  lexer.Token
	Name   *Identifier
	Params []*Identifier
	Body   *Block
}

func (f *Function) Literal() string { return f.Token.Literal }
func (f *Function) String() string {
	var out bytes.Buffer
	var params []string
	for _, p := range f.Params {
		params = append(params, p.String())
	}
	out.WriteString(f.Literal())
	out.WriteString(" ")
	out.WriteString(f.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	return out.String()
}

func (p *Parser) newFunction() Statement {
	fun := &Function{Token: p.token}
	if p.nextToken() != lexer.IDENT {
		err := util.NewError(p.token, util.ExpectedIdent, p.token.Literal)
		p.errors.Add(err)
		return nil
	}
	fun.Name = p.newIdentifier().(*Identifier)
	p.nextToken()
	if !p.isToken(lexer.LPAREN) {
		err := util.NewError(p.token, util.ExpectedParen, p.token.Literal)
		p.errors.Add(err)
		return nil
	}
	params := p.newParameters()
	for _, param := range params {
		ident, valid := param.(*Identifier)
		if !valid {
			err := util.NewError(p.token, util.ExpectedIdent, param.String())
			p.errors.Add(err)
			return nil
		}
		fun.Params = append(fun.Params, ident)
	}
	if p.nextToken() != lexer.LBRACE {
		err := util.NewError(p.token, util.ExpectedBrace, p.token.Literal)
		p.errors.Add(err)
		return nil
	}
	fun.Body = p.newBlock()
	return fun
}
