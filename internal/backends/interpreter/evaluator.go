package interpreter

import (
	"fmt"
	"sigil/internal/ast"
	"sigil/internal/backends"
)

type Evaluator struct{}

func NewEvaluator() backends.CompilerBackend {
	return &Evaluator{}
}

func (e *Evaluator) Execute(program *ast.Program, debug bool) error {
	var last Object
	for _, stmt := range program.Statements {
		val := Eval(stmt)

		if val != nil {
			last = val
		}
	}

	if debug {
		fmt.Printf("INTERPRET RESULT: %+v\n", last)
	}

	return nil
}

// Eval evaluates an AST Node and returns the resulting Object.
func Eval(node ast.Node) Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &ReturnObject{Value: val}

	// Expressions
	case *ast.NumberLiteral:
		return &Number{Value: node.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		return &String{Value: node.Value}

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}

		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}
