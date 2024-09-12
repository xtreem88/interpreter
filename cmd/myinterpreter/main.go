package main

import (
	"fmt"
	"math"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/astprinter"
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
	case "parse":
		parser := parser.NewParser(tokens)
		expression, err := parser.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
			os.Exit(65)
		}

		printer := astprinter.NewAstPrinter()
		fmt.Println(printer.Print(expression))
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	if scanner.HadError() {
		os.Exit(65)
	}
}
