package astprinter

import (
	"fmt"
	"math"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/parser"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(expr parser.Expr) string {
	return expr.Accept(a).(string)
}

func (a *AstPrinter) PrintStmt(stmt parser.Stmt) string {
	return stmt.Accept(a).(string)
}

func (a *AstPrinter) VisitBlockStmt(stmt *parser.BlockStmt) interface{} {
	return stmt.Accept(a)
}

func (a *AstPrinter) VisitLiteralExpr(expr *parser.Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	if number, ok := expr.Value.(float64); ok {
		if math.Floor(number) == number {
			return fmt.Sprintf("%.1f", number)
		}
		return fmt.Sprintf("%g", number)
	}
	if str, ok := expr.Value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *AstPrinter) VisitGroupingExpr(expr *parser.Grouping) interface{} {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitUnaryExpr(expr *parser.Unary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) VisitBinaryExpr(expr *parser.Binary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

// New method to implement
func (a *AstPrinter) VisitAssignExpr(expr *parser.Assign) interface{} {
	return a.parenthesize("assign "+expr.Name.Lexeme, expr.Value)
}

func (a *AstPrinter) VisitVariableExpr(expr *parser.Variable) interface{} {
	return expr.Name.Lexeme
}

// Existing visitor methods for statements
func (a *AstPrinter) VisitPrintStmt(stmt *parser.PrintStmt) interface{} {
	return a.parenthesize("print", stmt.Expression)
}

func (a *AstPrinter) VisitExpressionStmt(stmt *parser.ExpressionStmt) interface{} {
	return a.parenthesize(";", stmt.Expression)
}

func (a *AstPrinter) VisitVarStmt(stmt *parser.VarStmt) interface{} {
	if stmt.Initializer != nil {
		return a.parenthesize("var "+stmt.Name.Lexeme, stmt.Initializer)
	}
	return a.parenthesize("var " + stmt.Name.Lexeme)
}

func (a *AstPrinter) parenthesize(name string, exprs ...parser.Expr) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(a).(string))
	}
	builder.WriteString(")")
	return builder.String()
}
