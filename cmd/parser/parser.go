package parser

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
}

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

type ExprVisitor interface {
	VisitLiteralExpr(expr *Literal) interface{}
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (Expr, error) {
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(scanner.FALSE) {
		return &Literal{Value: false}, nil
	}
	if p.match(scanner.TRUE) {
		return &Literal{Value: true}, nil
	}
	if p.match(scanner.NIL) {
		return &Literal{Value: nil}, nil
	}
	if p.match(scanner.NUMBER) {
		return &Literal{Value: p.previous().Literal}, nil
	}
	if p.match(scanner.STRING) {
		return &Literal{Value: p.previous().Literal}, nil
	}

	return nil, fmt.Errorf("Expect expression.")
}

func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}
