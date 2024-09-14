package main

import (
	"fmt"
	"math"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/astprinter"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/interpreter"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/parser"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/scanner"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh [command] <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	scanner := scanner.NewScanner(string(fileContents))
	tokens := scanner.ScanTokens()

	switch command {
	case "tokenize":
		for _, token := range tokens {
			if token.Literal != nil {
				if number, ok := token.Literal.(float64); ok {
					if math.Floor(number) == number {
						fmt.Printf("%s %s %.1f\n", token.Type, token.Lexeme, number)
					} else {
						fmt.Printf("%s %s %g\n", token.Type, token.Lexeme, number)
					}
				} else {
					fmt.Printf("%s %s %v\n", token.Type, token.Lexeme, token.Literal)
				}
			} else {
				fmt.Printf("%s %s null\n", token.Type, token.Lexeme)
			}
		}

		if scanner.HadError() {
			os.Exit(65)
		}
	case "parse":
		if scanner.HadError() {
			os.Exit(65)
		}
		parser := parser.NewParser(tokens)
		expression, err := parser.ParseExpression()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
			os.Exit(65)
		}

		if expression == nil {
			fmt.Fprintln(os.Stderr, "No expression found to parse.")
			os.Exit(65)
		}

		printer := astprinter.NewAstPrinter()
		fmt.Println(printer.Print(expression))
	case "evaluate":
		if scanner.HadError() {
			os.Exit(65)
		}
		parser := parser.NewParser(tokens)
		expression, err := parser.ParseExpression()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
			os.Exit(70)
		}

		interpreter := interpreter.NewInterpreter()
		result, err := interpreter.Evaluate(expression)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(70)
		}
		fmt.Println(interpreter.Stringify(result))
	case "run":
		if scanner.HadError() {
			os.Exit(65)
		}
		parser := parser.NewParser(tokens)
		statements, err := parser.ParseStatements()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(65)
		}

		interpreter := interpreter.NewInterpreter()
		err = interpreter.Interpret(statements)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
			os.Exit(70)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
