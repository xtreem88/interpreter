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

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

type Unary struct {
	Operator scanner.Token
	Right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

type Binary struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

type PrintStmt struct {
	Expression Expr
}

func (p *PrintStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(p)
}

type ExpressionStmt struct {
	Expression Expr
}

func (e *ExpressionStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(e)
}

func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

type ExprVisitor interface {
	VisitLiteralExpr(expr *Literal) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitBinaryExpr(expr *Binary) interface{}
}

type StmtVisitor interface {
	VisitPrintStmt(stmt *PrintStmt) interface{}
	VisitExpressionStmt(stmt *ExpressionStmt) interface{}
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() (Expr, []Stmt, error) {
	var statements []Stmt
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, nil, err
		}
		statements = append(statements, stmt)
	}

	// For backwards compatibility, return the first expression if there's only one statement
	if len(statements) == 1 {
		if exprStmt, ok := statements[0].(*ExpressionStmt); ok {
			return exprStmt.Expression, statements, nil
		}
	}

	// If there are multiple statements or the single statement is not an expression,
	// return nil for the expression
	return nil, statements, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return &PrintStmt{Expression: value}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(scanner.SEMICOLON, "Expect ';' after expression.")

	return &ExpressionStmt{Expression: expr}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.EQUAL_EQUAL, scanner.BANG_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{Operator: operator, Right: right}, nil
	}

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
	if p.match(scanner.NUMBER, scanner.STRING) {
		return &Literal{Value: p.previous().Literal}, nil
	}
	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &Grouping{Expression: expr}, nil
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

func (p *Parser) consume(t scanner.TokenType, message string) (scanner.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return scanner.Token{}, fmt.Errorf("%s", message)
}
