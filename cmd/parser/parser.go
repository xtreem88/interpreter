package parser

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
	errors  []error
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
	VisitVariableExpr(expr *Variable) interface{}
	VisitAssignExpr(expr *Assign) interface{}
}

type StmtVisitor interface {
	VisitPrintStmt(stmt *PrintStmt) interface{}
	VisitExpressionStmt(stmt *ExpressionStmt) interface{}
	VisitVarStmt(stmt *VarStmt) interface{}
}

type Variable struct {
	Name scanner.Token
}

func (v *Variable) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(v)
}

type Assign struct {
	Name  scanner.Token
	Value Expr
}

func (a *Assign) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(a)
}

type VarStmt struct {
	Name        scanner.Token
	Initializer Expr
}

func (v *VarStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(v)
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() (Expr, []Stmt, error) {
	var statements []Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			p.synchronize()
			continue
		}
		statements = append(statements, stmt)
	}

	// If there is exactly one statement and it's an expression statement,
	// return its expression for the 'evaluate' command.
	if len(statements) == 1 {
		if exprStmt, ok := statements[0].(*ExpressionStmt); ok {
			return exprStmt.Expression, statements, nil
		} else if printStmt, ok := statements[0].(*PrintStmt); ok {
			// For 'print' statements, we can handle them as expressions
			return printStmt.Expression, statements, nil
		}
	}

	// Otherwise, return nil for the expression
	return nil, statements, nil
}

func (p *Parser) ParseExpression() (Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.isAtEnd() {
		return nil, p.error(p.peek(), "Unexpected tokens after expression.")
	}
	return expr, nil
}

func (p *Parser) ParseStatements() ([]Stmt, error) {
	var statements []Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			p.synchronize()
			continue
		}
		statements = append(statements, stmt)
	}

	if len(p.errors) > 0 {
		errorMessages := ""
		for _, err := range p.errors {
			errorMessages += err.Error() + "\n"
		}
		return statements, fmt.Errorf(errorMessages)
	}

	return statements, nil
}

func (p *Parser) error(token scanner.Token, message string) error {
	var errMsg string
	if token.Type == scanner.EOF {
		errMsg = fmt.Sprintf("[line %d] Error at end: %s", token.Line, message)
	} else {
		errMsg = fmt.Sprintf("[line %d] Error at '%s': %s", token.Line, token.Lexeme, message)
	}
	err := fmt.Errorf(errMsg)
	p.errors = append(p.errors, err)

	// fmt.Fprintln(os.Stderr, errMsg)
	return err
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(scanner.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(scanner.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return &VarStmt{Name: name, Initializer: initializer}, nil
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == scanner.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case scanner.CLASS, scanner.FUN, scanner.VAR, scanner.FOR, scanner.IF, scanner.WHILE, scanner.PRINT, scanner.RETURN:
			return
		}

		p.advance()
	}
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
	if value == nil {
		return nil, p.error(p.peek(), "Expect expression after 'print'.")
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
	if expr == nil {
		return nil, p.error(p.peek(), "Expect expression.")
	}
	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}
	return &ExpressionStmt{Expression: expr}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
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

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(scanner.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if variable, ok := expr.(*Variable); ok {
			return &Assign{Name: variable.Name, Value: value}, nil
		}

		return nil, p.error(equals, "Invalid assignment target.")
	}

	return expr, nil
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
	if p.match(scanner.IDENTIFIER) {
		return &Variable{Name: p.previous()}, nil
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

	return nil, p.error(p.peek(), "Expect expression.")
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
