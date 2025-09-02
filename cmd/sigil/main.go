package main

import (
	"fmt"
	"os"
	"sigil/internal/backends/interpreter"
	"sigil/internal/lexer"
	"sigil/internal/parser"
	"sigil/internal/typechecker"
)

const DEBUG_MODE = true

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Sigil Language Compiler")
		fmt.Println("Usage: sigil <file.sgl>")
		return
	}

	filename := os.Args[1]

	// Read the source file
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		return
	}

	if DEBUG_MODE {

		fmt.Printf("Compiling: %s\n", filename)
		fmt.Printf("Source code:\n%s\n", string(source))
		fmt.Println("\nTokens:")

		// Tokenize the source
		l := lexer.New(string(source))
		for {
			tok := l.NextToken()
			fmt.Printf("Type: %v, Literal: '%s', Line: %d, Column: %d\n",
				tok.Type, tok.Literal, tok.Line, tok.Column)
			if tok.Type == lexer.EOF {
				break
			}
		}

		fmt.Println("\nParsing:")

	}

	// Reset lexer for parsing
	l := lexer.New(string(source))
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parser errors
	if errors := p.Errors(); len(errors) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range errors {
			fmt.Printf(" %s\n", err)
		}
		return // Don't continue if there are parse errors
	}

	if DEBUG_MODE {
		fmt.Printf("AST Tree:\n%s", program.TreeString("", false))
		fmt.Println("\nType Checking:")
	}

	// Type check the program
	tc := typechecker.New()
	tc.CheckProgram(program)

	if tc.HasErrors() {
		fmt.Println("Type errors:")
		for _, err := range tc.Errors() {
			fmt.Printf(" %s\n", err.Error())
		}
		return // Don't continue if there are type errors
	}

	if DEBUG_MODE {
		fmt.Println("✓ Type checking passed")
		fmt.Println("\nExecution:")
	}

	// Execute the program with the interpreter backend
	backend := interpreter.NewEvaluator()
	err = backend.Execute(program, DEBUG_MODE)

	if err != nil {
		fmt.Printf("Runtime error: %s\n", err)
		return
	}
}
