package astprinter

import (
	"fmt"
	"math"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/parser"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(expr parser.Expr) string {
	return expr.Accept(a).(string)
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
	return fmt.Sprintf("%v", expr.Value)
}
