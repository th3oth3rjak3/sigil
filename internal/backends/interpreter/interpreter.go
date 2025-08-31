package interpreter

import (
	"fmt"
	"sigil/internal/ast"
	"sigil/internal/backends"
	"strconv"
)

// Value represents a runtime value in the interpreter
type Value interface {
	String() string
	Type() string
}

// Runtime value types
type NumberValue struct {
	Value float64
}

func (nv *NumberValue) String() string { return fmt.Sprintf("%g", nv.Value) }
func (nv *NumberValue) Type() string   { return "Number" }

type StringValue struct {
	Value string
}

func (sv *StringValue) String() string { return sv.Value }
func (sv *StringValue) Type() string   { return "String" }

type BoolValue struct {
	Value bool
}

func (bv *BoolValue) String() string { return fmt.Sprintf("%t", bv.Value) }
func (bv *BoolValue) Type() string   { return "Bool" }

// Runtime environment for variable storage
type Environment struct {
	store map[string]Value
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Value),
		outer: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Value, bool) {
	value, ok := e.store[name]
	if !ok && e.outer != nil {
		value, ok = e.outer.Get(name)
	}
	return value, ok
}

func (e *Environment) Set(name string, value Value) {
	e.store[name] = value
}

// Interpreter implements the CompilerBackend interface
type Interpreter struct {
	env *Environment
}

// New creates a new interpreter instance
func New() backends.CompilerBackend {
	return &Interpreter{
		env: NewEnvironment(),
	}
}

// Execute implements the CompilerBackend interface
func (i *Interpreter) Execute(program *ast.Program) error {
	for _, stmt := range program.Statements {
		err := i.executeStatement(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) executeStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.LetStatement:
		return i.executeLetStatement(s)
	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (i *Interpreter) executeLetStatement(stmt *ast.LetStatement) error {
	value, err := i.evaluateExpression(stmt.Value)
	if err != nil {
		return err
	}
	i.env.Set(stmt.Name.Value, value)
	return nil
}

func (i *Interpreter) evaluateExpression(expr ast.Expression) (Value, error) {
	switch e := expr.(type) {
	case *ast.NumberLiteral:
		return i.evaluateNumberLiteral(e)
	case *ast.StringLiteral:
		return i.evaluateStringLiteral(e)
	case *ast.BooleanLiteral:
		return i.evaluateBooleanLiteral(e)
	case *ast.Identifier:
		return i.evaluateIdentifier(e)
	case *ast.InfixExpression:
		return i.evaluateInfixExpression(e)
	case *ast.PrefixExpression:
		return i.evaluatePrefixExpression(e)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (i *Interpreter) evaluateNumberLiteral(literal *ast.NumberLiteral) (Value, error) {
	value, err := strconv.ParseFloat(literal.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number literal: %s", literal.Value)
	}
	return &NumberValue{Value: value}, nil
}

func (i *Interpreter) evaluateStringLiteral(literal *ast.StringLiteral) (Value, error) {
	return &StringValue{Value: literal.String()}, nil
}

func (i *Interpreter) evaluateBooleanLiteral(literal *ast.BooleanLiteral) (Value, error) {
	return &BoolValue{Value: literal.Value}, nil
}

func (i *Interpreter) evaluateIdentifier(ident *ast.Identifier) (Value, error) {
	value, exists := i.env.Get(ident.Value)
	if !exists {
		return nil, fmt.Errorf("undefined variable: %s", ident.Value)
	}
	return value, nil
}

func (i *Interpreter) evaluateInfixExpression(expr *ast.InfixExpression) (Value, error) {
	left, err := i.evaluateExpression(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluateExpression(expr.Right)
	if err != nil {
		return nil, err
	}

	return i.applyInfixOperator(expr.Operator, left, right)
}

func (i *Interpreter) evaluatePrefixExpression(expr *ast.PrefixExpression) (Value, error) {
	operand, err := i.evaluateExpression(expr.Right)
	if err != nil {
		return nil, err
	}

	return i.applyPrefixOperator(expr.Operator, operand)
}

func (i *Interpreter) applyInfixOperator(operator string, left, right Value) (Value, error) {
	switch operator {
	case "+", "-", "*", "/":
		return i.applyArithmeticOperator(operator, left, right)
	case "==", "!=":
		return i.applyEqualityOperator(operator, left, right)
	case "<", ">", "<=", ">=":
		return i.applyComparisonOperator(operator, left, right)
	default:
		return nil, fmt.Errorf("unknown infix operator: %s", operator)
	}
}

func (i *Interpreter) applyArithmeticOperator(operator string, left, right Value) (Value, error) {
	leftNum, leftOk := left.(*NumberValue)
	rightNum, rightOk := right.(*NumberValue)

	if !leftOk || !rightOk {
		return nil, fmt.Errorf("arithmetic operators require numbers")
	}

	switch operator {
	case "+":
		return &NumberValue{Value: leftNum.Value + rightNum.Value}, nil
	case "-":
		return &NumberValue{Value: leftNum.Value - rightNum.Value}, nil
	case "*":
		return &NumberValue{Value: leftNum.Value * rightNum.Value}, nil
	case "/":
		if rightNum.Value == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return &NumberValue{Value: leftNum.Value / rightNum.Value}, nil
	default:
		return nil, fmt.Errorf("unknown arithmetic operator: %s", operator)
	}
}

func (i *Interpreter) applyEqualityOperator(operator string, left, right Value) (Value, error) {
	switch operator {
	case "==":
		return &BoolValue{Value: i.valuesEqual(left, right)}, nil
	case "!=":
		return &BoolValue{Value: !i.valuesEqual(left, right)}, nil
	default:
		return nil, fmt.Errorf("unknown equality operator: %s", operator)
	}
}

func (i *Interpreter) applyComparisonOperator(operator string, left, right Value) (Value, error) {
	leftNum, leftOk := left.(*NumberValue)
	rightNum, rightOk := right.(*NumberValue)

	if !leftOk || !rightOk {
		return nil, fmt.Errorf("comparison operators require numbers")
	}

	switch operator {
	case "<":
		return &BoolValue{Value: leftNum.Value < rightNum.Value}, nil
	case ">":
		return &BoolValue{Value: leftNum.Value > rightNum.Value}, nil
	case "<=":
		return &BoolValue{Value: leftNum.Value <= rightNum.Value}, nil
	case ">=":
		return &BoolValue{Value: leftNum.Value >= rightNum.Value}, nil
	default:
		return nil, fmt.Errorf("unknown comparison operator: %s", operator)
	}
}

func (i *Interpreter) applyPrefixOperator(operator string, operand Value) (Value, error) {
	switch operator {
	case "-":
		num, ok := operand.(*NumberValue)
		if !ok {
			return nil, fmt.Errorf("unary minus requires a number")
		}
		return &NumberValue{Value: -num.Value}, nil
	case "!":
		bool, ok := operand.(*BoolValue)
		if !ok {
			return nil, fmt.Errorf("logical not requires a boolean")
		}
		return &BoolValue{Value: !bool.Value}, nil
	default:
		return nil, fmt.Errorf("unknown prefix operator: %s", operator)
	}
}

func (i *Interpreter) valuesEqual(left, right Value) bool {
	if left.Type() != right.Type() {
		return false
	}

	switch l := left.(type) {
	case *NumberValue:
		r := right.(*NumberValue)
		return l.Value == r.Value
	case *StringValue:
		r := right.(*StringValue)
		return l.Value == r.Value
	case *BoolValue:
		r := right.(*BoolValue)
		return l.Value == r.Value
	default:
		return false
	}
}
