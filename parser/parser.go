package parser

import (
	"Eldrlang/lexer"
	"Eldrlang/util"
)

type Parser struct {
	lexer  *lexer.Lexer
	token  lexer.Token
	errors util.Errors
}

func NewParser() *Parser {
	return &Parser{lexer: lexer.NewLexer(nil)}
}

func (p *Parser) Errors() util.Errors {
	return p.errors
}

func (p *Parser) ParsePackage(input, pkgName string) *Package {
	var (
		pkg  = Package{Namespace: pkgName}
		node Node
	)
	p.lexer.UpdateInput([]byte(input))
	for p.nextToken() != lexer.EOF {
		node = p.parseStatement()
		if node != nil { //TODO FIX NULL RETURNS AT { 5 + (a * 2) };
			pkg.Nodes = append(pkg.Nodes, node)
		}
	}
	pkg.Errors = p.errors
	p.errors.Clear()
	return &pkg
}

func (p *Parser) parseStatement() Statement {
	switch p.token.Type {
	case lexer.VARIABLE:
		return p.newVariable()
	case lexer.LBRACE:
		return p.newBlock()
	case lexer.IF:
		return p.newConditional()
	case lexer.LOOP:
		return p.newLoop()
	case lexer.FUNCTION:
		return p.newFunction()
	case lexer.RETURN:
		return p.newReturn()
	case lexer.BREAK:
		return p.newBreak()
	default:
		return p.parseExpression()
	}
}

func (p *Parser) parseExpression() Expression {
	var exp Expression
	for !p.isTokenOrEOF(lexer.SEMICOLON) {
		exp = p.parseToken(exp)
		p.nextToken()
	}
	if p.isToken(lexer.EOF) {
		err := util.NewError(p.token, util.UnexpectedEOF, ";")
		p.errors.Add(err)
		return nil
	}
	return exp
}

func (p *Parser) parseGroupExpression() Expression {
	p.nextToken() // Skip ( opening
	var exp Expression
	for !p.isTokenOrEOF(lexer.RPAREN) {
		exp = p.parseToken(exp)
		p.nextToken()
	}
	if p.isToken(lexer.EOF) {
		err := util.NewError(p.token, util.UnexpectedEOF, ")")
		p.errors.Add(err)
		return nil
	}
	return exp
}

func (p *Parser) parseToken(exp Expression) Expression {
	switch p.token.Type {
	case lexer.IDENT:
		exp = p.newIdentifier()
		if p.isPeekToken(lexer.LPAREN) {
			p.nextToken() // Skip current ident token
			exp = p.newFuncCall(exp)
		}
	case lexer.TRUE, lexer.FALSE:
		exp = p.newBoolean()
	case lexer.INTEGER:
		exp = p.newInteger()
	case lexer.STRING:
		exp = p.newString()
	case lexer.BANG:
		exp = p.newPrefix()
	case lexer.PLUS, lexer.MINUS:
		switch exp.(type) {
		case *String, nil:
			exp = p.newPrefix()
		default:
			exp = p.newInfix(exp)
		}
	case lexer.ASSIGN:
		switch exp.(type) {
		case *Variable, *Identifier:
			exp = p.newAssign(exp)
		default:
			err := util.NewError(p.token, util.IllegalOpeAtt)
			p.errors.Add(err)
			return nil
		}
	case lexer.EQ, lexer.NOTEQ, lexer.LT:
		fallthrough
	case lexer.LTEQ, lexer.GT, lexer.GTEQ:
		fallthrough
	case lexer.ASTERISK, lexer.SLASH:
		exp = p.newInfix(exp)
	case lexer.LPAREN:
		exp = p.parseGroupExpression()
	default:
		err := util.NewError(p.token, util.IllegalLetter, p.token.Literal)
		p.errors.Add(err)
		return nil
	}
	return exp
}

func (p *Parser) nextToken() lexer.Type {
	p.token = p.lexer.NextToken()
	return p.token.Type
}

func (p *Parser) peekToken() lexer.Type {
	return p.lexer.PeekToken().Type
}

func (p *Parser) isToken(t lexer.Type) bool {
	return p.token.Type == t
}

func (p *Parser) isPeekToken(t lexer.Type) bool {
	return p.peekToken() == t
}

func (p *Parser) isTokenOrEOF(t lexer.Type) bool {
	return p.isToken(t) || p.isToken(lexer.EOF)
}
