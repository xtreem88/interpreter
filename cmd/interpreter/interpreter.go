package interpreter

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/parser"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/scanner"
)

type RuntimeError struct {
	token   scanner.Token
	message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", e.message, e.token.Line)
}

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{environment: NewEnvironment()}
}

func (i *Interpreter) Evaluate(expr parser.Expr) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(RuntimeError); ok {
				err = &runtimeErr
			} else {
				panic(r)
			}
		}
	}()

	result = expr.Accept(i)
	return result, nil
}

func (i *Interpreter) Interpret(statements []parser.Stmt) error {
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) execute(stmt parser.Stmt) error {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*RuntimeError); ok {
				fmt.Fprintln(os.Stderr, runtimeErr.Error())
			} else {
				panic(r)
			}
		}
	}()

	result := stmt.Accept(i)
	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *parser.PrintStmt) interface{} {
	value, err := i.Evaluate(stmt.Expression)
	if err != nil {
		return err
	}
	fmt.Println(i.stringify(value))
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *parser.ExpressionStmt) interface{} {
	_, err := i.Evaluate(stmt.Expression)
	return err
}

func (i *Interpreter) stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}
	if b, ok := value.(bool); ok {
		return strconv.FormatBool(b)
	}
	if n, ok := value.(float64); ok {
		return strconv.FormatFloat(n, 'f', -1, 64)
	}
	if s, ok := value.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", value)
}

func (i *Interpreter) VisitLiteralExpr(expr *parser.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *parser.Grouping) interface{} {
	eval, _ := i.Evaluate(expr.Expression)
	return eval
}

func (i *Interpreter) VisitUnaryExpr(expr *parser.Unary) interface{} {
	right := expr.Right.Accept(i)

	switch expr.Operator.Type {
	case scanner.MINUS:
		if num, ok := right.(float64); ok {
			return -num
		}
		panic(RuntimeError{token: expr.Operator, message: "Operand must be a number."})
	case scanner.BANG:
		return !i.isTruthy(right)
	}

	// Unreachable
	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr *parser.Binary) interface{} {
	left := expr.Left.Accept(i)
	right := expr.Right.Accept(i)

	switch expr.Operator.Type {
	case scanner.MINUS:
		return i.checkNumberOperands(expr.Operator, left, right)
	case scanner.PLUS:
		if leftStr, leftOk := left.(string); leftOk {
			if rightStr, rightOk := right.(string); rightOk {
				return leftStr + rightStr
			}
		}
		return i.checkNumberOperands(expr.Operator, left, right)
	case scanner.STAR:
		return i.checkNumberOperands(expr.Operator, left, right)
	case scanner.SLASH:
		return i.checkNumberOperands(expr.Operator, left, right)
	case scanner.GREATER:
		return i.checkNumberOperands(expr.Operator, left, right)
	case scanner.GREATER_EQUAL:
		return i.checkNumberOperands(expr.Operator, left, right)
	case scanner.LESS:
		return i.checkNumberOperands(expr.Operator, left, right)
	case scanner.LESS_EQUAL:
		return i.checkNumberOperands(expr.Operator, left, right)
	case scanner.EQUAL_EQUAL:
		return i.isEqual(left, right)
	case scanner.BANG_EQUAL:
		return !i.isEqual(left, right)
	}

	// Unreachable
	return nil
}

func (i *Interpreter) VisitVariableExpr(expr *parser.Variable) interface{} {
	value, err := i.environment.Get(expr.Name)
	if err != nil {
		panic(err)
	}
	return value
}

func (i *Interpreter) VisitAssignExpr(expr *parser.Assign) interface{} {
	value := expr.Value.Accept(i)

	err := i.environment.Assign(expr.Name, value)
	if err != nil {
		panic(err)
	}
	return value
}

func (i *Interpreter) VisitVarStmt(stmt *parser.VarStmt) interface{} {
	var value interface{}
	if stmt.Initializer != nil {
		value = stmt.Initializer.Accept(i)
	}

	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) checkNumberOperands(operator scanner.Token, left, right interface{}) interface{} {
	leftNum, leftOk := left.(float64)
	rightNum, rightOk := right.(float64)
	if !leftOk || !rightOk {
		panic(RuntimeError{token: operator, message: "Operands must be numbers."})
	}

	switch operator.Type {
	case scanner.MINUS:
		return leftNum - rightNum
	case scanner.PLUS:
		return leftNum + rightNum
	case scanner.STAR:
		return leftNum * rightNum
	case scanner.SLASH:
		if rightNum == 0 {
			panic(RuntimeError{token: operator, message: "Division by zero."})
		}
		return leftNum / rightNum
	case scanner.GREATER:
		return leftNum > rightNum
	case scanner.GREATER_EQUAL:
		return leftNum >= rightNum
	case scanner.LESS:
		return leftNum < rightNum
	case scanner.LESS_EQUAL:
		return leftNum <= rightNum
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
