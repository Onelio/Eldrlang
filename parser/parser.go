package parser

import (
	"Eldrlang/lexer"
	"Eldrlang/util"
)

type Parser struct {
	errors util.Errors
	lexer  *lexer.Lexer
	token  lexer.Token
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseProgram(input string) *Program {
	var (
		program Program
		node    Node
	)
	p.lexer = lexer.NewLexer([]byte(input))
	p.errors = util.Errors{}
	for p.nextToken() != lexer.EOF {
		node = p.parseStatement()
		if node != nil {
			program.Nodes = append(program.Nodes, node)
		}
	}
	return &program
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
		switch exp.(type) {
		case *Identifier:
			exp = p.newFuncCall(exp)
		default:
			exp = p.parseGroupExpression()
		}
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
