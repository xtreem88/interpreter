package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/parser"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/scanner"
)

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Evaluate(expr parser.Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(expr *parser.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *parser.Grouping) interface{} {
	return i.Evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *parser.Unary) interface{} {
	right := i.Evaluate(expr.Right)

	switch expr.Operator.Type {
	case scanner.MINUS:
		if num, ok := right.(float64); ok {
			return -num
		}
		// Handle error: Operand must be a number
		panic(fmt.Sprintf("Operand must be a number"))
	case scanner.BANG:
		return !i.isTruthy(right)
	}

	// Unreachable
	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr *parser.Binary) interface{} {
	// We'll implement this later
	return nil
}

func (i *Interpreter) Stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}
	if b, ok := object.(bool); ok {
		if b {
			return "true"
		}
		return "false"
	}
	if num, ok := object.(float64); ok {
		text := fmt.Sprintf("%g", num)
		if text[len(text)-2:] == ".0" {
			text = text[:len(text)-2]
		}
		return text
	}
	return fmt.Sprintf("%v", object)
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(bool); ok {
		return b
	}
	return true
}
