package interpreter

import (
	"fmt"
	"strconv"

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
		return -i.checkNumberOperand(right)
	case scanner.BANG:
		return !i.isTruthy(right)
	}

	// Unreachable
	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr *parser.Binary) interface{} {
	left := i.Evaluate(expr.Left)
	right := i.Evaluate(expr.Right)

	switch expr.Operator.Type {
	case scanner.MINUS:
		return i.checkNumberOperand(left) - i.checkNumberOperand(right)
	case scanner.PLUS:
		if leftStr, leftOk := left.(string); leftOk {
			if rightStr, rightOk := right.(string); rightOk {
				return leftStr + rightStr
			}
		}
		return i.checkNumberOperand(left) + i.checkNumberOperand(right)
	case scanner.STAR:
		return i.checkNumberOperand(left) * i.checkNumberOperand(right)
	case scanner.SLASH:
		rightNum := i.checkNumberOperand(right)
		if rightNum == 0 {
			panic(fmt.Sprintf("Division by zero"))
		}
		return i.checkNumberOperand(left) / rightNum
	case scanner.GREATER:
		return i.checkNumberOperand(left) > i.checkNumberOperand(right)
	case scanner.GREATER_EQUAL:
		return i.checkNumberOperand(left) >= i.checkNumberOperand(right)
	case scanner.LESS:
		return i.checkNumberOperand(left) < i.checkNumberOperand(right)
	case scanner.LESS_EQUAL:
		return i.checkNumberOperand(left) <= i.checkNumberOperand(right)
	case scanner.EQUAL_EQUAL:
		return i.isEqual(left, right)
	case scanner.BANG_EQUAL:
		return !i.isEqual(left, right)
	}

	// Unreachable
	return nil
}

func (i *Interpreter) isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}

func (i *Interpreter) checkNumberOperand(operand interface{}) float64 {
	if num, ok := operand.(float64); ok {
		return num
	}
	panic(fmt.Sprintf("Operand must be a number"))
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
		return strconv.FormatFloat(num, 'f', -1, 64)
	}
	if str, ok := object.(string); ok {
		return str
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
