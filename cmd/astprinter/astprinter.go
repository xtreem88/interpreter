package astprinter

import (
	"fmt"

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
	return fmt.Sprintf("%v", expr.Value)
}
