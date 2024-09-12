package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/parser"
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
	// We'll implement this later
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
	return fmt.Sprintf("%v", object)
}
