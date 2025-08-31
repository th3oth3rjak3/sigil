package typechecker

import (
	"fmt"
	"sigil/internal/ast"
)

// Type represents a type in the language
type Type interface {
	String() string
	Equals(Type) bool
}

// Basic types
type NumberType struct{}
type StringType struct{}
type BoolType struct{}
type VoidType struct{} // for statements that don't return values

func (nt *NumberType) String() string { return "Number" }
func (nt *NumberType) Equals(other Type) bool {
	_, ok := other.(*NumberType)
	return ok
}

func (st *StringType) String() string { return "String" }
func (st *StringType) Equals(other Type) bool {
	_, ok := other.(*StringType)
	return ok
}

func (bt *BoolType) String() string { return "Bool" }
func (bt *BoolType) Equals(other Type) bool {
	_, ok := other.(*BoolType)
	return ok
}

func (vt *VoidType) String() string { return "Void" }
func (vt *VoidType) Equals(other Type) bool {
	_, ok := other.(*VoidType)
	return ok
}

// Type error with position information
type TypeError struct {
	Message string
	Line    int
	Column  int
}

func (te *TypeError) Error() string {
	return fmt.Sprintf("Type error at line %d, column %d: %s", te.Line, te.Column, te.Message)
}

// Symbol represents a variable in the symbol table
type Symbol struct {
	Name   string
	Type   Type
	Line   int
	Column int
}

// Environment for symbol table with scope support
type Environment struct {
	store map[string]*Symbol
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]*Symbol),
		outer: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (*Symbol, bool) {
	symbol, ok := e.store[name]
	if !ok && e.outer != nil {
		symbol, ok = e.outer.Get(name)
	}
	return symbol, ok
}

func (e *Environment) Set(name string, symbol *Symbol) {
	e.store[name] = symbol
}

// Type Checker
type TypeChecker struct {
	env    *Environment
	errors []*TypeError
}

func New() *TypeChecker {
	return &TypeChecker{
		env:    NewEnvironment(),
		errors: []*TypeError{},
	}
}

func (tc *TypeChecker) addError(message string, line, column int) {
	tc.errors = append(tc.errors, &TypeError{
		Message: message,
		Line:    line,
		Column:  column,
	})
}

func (tc *TypeChecker) Errors() []*TypeError {
	return tc.errors
}

func (tc *TypeChecker) HasErrors() bool {
	return len(tc.errors) > 0
}

// Helper function to parse type annotations from identifiers
func (tc *TypeChecker) parseTypeFromIdentifier(typeIdent *ast.Identifier) Type {
	if typeIdent == nil {
		return nil
	}

	switch typeIdent.Value {
	case "Number":
		return &NumberType{}
	case "String":
		return &StringType{}
	case "Bool":
		return &BoolType{}
	default:
		tc.addError(fmt.Sprintf("unknown type: %s", typeIdent.Value), typeIdent.Token.Line, typeIdent.Token.Column)
		return nil
	}
}

// Main entry point for type checking
func (tc *TypeChecker) CheckProgram(program *ast.Program) Type {
	for _, stmt := range program.Statements {
		tc.CheckStatement(stmt)
	}
	return &VoidType{}
}

func (tc *TypeChecker) CheckStatement(stmt ast.Statement) Type {
	switch s := stmt.(type) {
	case *ast.PrintStatement:
		return tc.CheckPrintStatement(s)
	case *ast.LetStatement:
		return tc.CheckLetStatement(s)
	default:
		tc.addError(fmt.Sprintf("unknown statement type: %T", stmt), 0, 0)
		return &VoidType{}
	}
}

func (tc *TypeChecker) CheckPrintStatement(stmt *ast.PrintStatement) Type {
	// Print statements can print any type, so just check the expression is valid
	tc.CheckExpression(stmt.Expression)
	return &VoidType{}
}

func (tc *TypeChecker) CheckLetStatement(stmt *ast.LetStatement) Type {
	// Type hint is mandatory for now
	if stmt.TypeHint == nil {
		tc.addError("type annotation is required for variable declarations", stmt.Token.Line, stmt.Token.Column)
		return &VoidType{}
	}

	// Parse the declared type
	declaredType := tc.parseTypeFromIdentifier(stmt.TypeHint)
	if declaredType == nil {
		return &VoidType{} // Error already reported in parseTypeFromIdentifier
	}

	// Check the value expression
	valueType := tc.CheckExpression(stmt.Value)

	// Verify the value type matches the declared type
	if valueType != nil && !declaredType.Equals(valueType) {
		tc.addError(
			fmt.Sprintf("type mismatch: declared %s but got %s", declaredType.String(), valueType.String()),
			stmt.Token.Line, stmt.Token.Column,
		)
	}

	// Add the variable to the environment
	symbol := &Symbol{
		Name:   stmt.Name.Value,
		Type:   declaredType,
		Line:   stmt.Name.Token.Line,
		Column: stmt.Name.Token.Column,
	}
	tc.env.Set(stmt.Name.Value, symbol)

	return &VoidType{}
}

func (tc *TypeChecker) CheckExpression(expr ast.Expression) Type {
	switch e := expr.(type) {
	case *ast.NumberLiteral:
		return &NumberType{}
	case *ast.BooleanLiteral:
		return &BoolType{}
	case *ast.Identifier:
		return tc.CheckIdentifier(e)
	case *ast.InfixExpression:
		return tc.CheckInfixExpression(e)
	case *ast.PrefixExpression:
		return tc.CheckPrefixExpression(e)
	default:
		tc.addError(fmt.Sprintf("unknown expression type: %T", expr), 0, 0)
		return &VoidType{}
	}
}

func (tc *TypeChecker) CheckIdentifier(ident *ast.Identifier) Type {
	symbol, exists := tc.env.Get(ident.Value)
	if !exists {
		tc.addError(fmt.Sprintf("undefined variable: %s", ident.Value), ident.Token.Line, ident.Token.Column)
		return &VoidType{}
	}
	return symbol.Type
}

func (tc *TypeChecker) CheckInfixExpression(expr *ast.InfixExpression) Type {
	leftType := tc.CheckExpression(expr.Left)
	rightType := tc.CheckExpression(expr.Right)

	switch expr.Operator {
	case "+", "-", "*", "/":
		// Arithmetic operators require numbers
		if !leftType.Equals(&NumberType{}) {
			tc.addError(fmt.Sprintf("left operand of %s must be Number, got %s", expr.Operator, leftType.String()), expr.Token.Line, expr.Token.Column)
		}
		if !rightType.Equals(&NumberType{}) {
			tc.addError(fmt.Sprintf("right operand of %s must be Number, got %s", expr.Operator, rightType.String()), expr.Token.Line, expr.Token.Column)
		}
		return &NumberType{}

	case "==", "!=":
		// Equality operators require same types
		if !leftType.Equals(rightType) {
			tc.addError(fmt.Sprintf("cannot compare %s with %s", leftType.String(), rightType.String()), expr.Token.Line, expr.Token.Column)
		}
		return &BoolType{}

	case "<", ">", "<=", ">=":
		// Comparison operators require numbers
		if !leftType.Equals(&NumberType{}) {
			tc.addError(fmt.Sprintf("left operand of %s must be Number, got %s", expr.Operator, leftType.String()), expr.Token.Line, expr.Token.Column)
		}
		if !rightType.Equals(&NumberType{}) {
			tc.addError(fmt.Sprintf("right operand of %s must be Number, got %s", expr.Operator, rightType.String()), expr.Token.Line, expr.Token.Column)
		}
		return &BoolType{}

	default:
		tc.addError(fmt.Sprintf("unknown infix operator: %s", expr.Operator), expr.Token.Line, expr.Token.Column)
		return &VoidType{}
	}
}

func (tc *TypeChecker) CheckPrefixExpression(expr *ast.PrefixExpression) Type {
	operandType := tc.CheckExpression(expr.Right)

	switch expr.Operator {
	case "-":
		if !operandType.Equals(&NumberType{}) {
			tc.addError(fmt.Sprintf("unary minus requires Number, got %s", operandType.String()), expr.Token.Line, expr.Token.Column)
		}
		return &NumberType{}

	case "!":
		if !operandType.Equals(&BoolType{}) {
			tc.addError(fmt.Sprintf("logical not requires Bool, got %s", operandType.String()), expr.Token.Line, expr.Token.Column)
		}
		return &BoolType{}

	default:
		tc.addError(fmt.Sprintf("unknown prefix operator: %s", expr.Operator), expr.Token.Line, expr.Token.Column)
		return &VoidType{}
	}
}
