package parser

import (
	"Eldrlang/lexer"
	"Eldrlang/util"
	"bytes"
	"strings"
)

type Variable struct {
	Token lexer.Token
	Name  *Identifier
}

func (d *Variable) Literal() string { return d.Token.Literal }
func (d *Variable) String() string {
	var out bytes.Buffer
	out.WriteString(d.Literal() + " ")
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

	p.nextToken()
	return p.parseToken(def)
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
	Init  Expression
	Cond  Expression
	Iter  Expression
	Body  *Block
}

func (l *Loop) Literal() string { return l.Token.Literal }
func (l *Loop) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	if l.Init != nil {
		out.WriteString("(")
		out.WriteString(l.Init.String())
		out.WriteString("; ")
		out.WriteString(l.Cond.String())
		out.WriteString("; ")
		out.WriteString(l.Iter.String())
		out.WriteString(";) ")
	}
	out.WriteString(l.Body.String())
	return out.String()
}

func (p *Parser) newLoop() Statement {
	loop := &Loop{Token: p.token, Cond: &Boolean{Value: true}}
	if p.nextToken() == lexer.LPAREN {
		p.nextToken() // Skip ( opening
		loop.Init = p.parseStatement()
		if _, valid := loop.Init.(*Assign); !valid {
			err := util.NewError(p.token, util.ExpectedInitL)
			p.errors.Add(err)
			return nil
		}
		p.nextToken() // Skip ;
		loop.Cond = p.parseExpression()
		if _, valid := loop.Cond.(*Infix); !valid {
			err := util.NewError(p.token, util.ExpectedCondL)
			p.errors.Add(err)
			return nil
		}
		p.nextToken() // Skip ;
		loop.Iter = p.parseExpression()
		p.nextToken() // Skip ;
		if !p.isToken(lexer.RPAREN) {
			err := util.NewError(p.token, util.ExpectedLoopC)
			p.errors.Add(err)
			return nil
		}
		p.nextToken() // Skip )
	}

	if !p.isToken(lexer.LBRACE) {
		err := util.NewError(p.token, util.ExpectedBrace, p.token.Literal)
		p.errors.Add(err)
		return nil
	}
	loop.Body = p.newBlock()
	return loop
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
	if !p.isPeekToken(lexer.LPAREN) {
		err := util.NewError(p.token, util.ExpectedParen, p.token.Literal)
		p.errors.Add(err)
		return nil
	}
	fun.Params = p.newParameters()
	if p.nextToken() != lexer.LBRACE {
		err := util.NewError(p.token, util.ExpectedBrace, p.token.Literal)
		p.errors.Add(err)
		return nil
	}
	fun.Body = p.newBlock()
	return fun
}
