package typechecker

import (
	"fmt"
	"sigil/internal/ast"
)

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
	// Look in the current environment
	symbol, ok := e.store[name]
	if ok {
		return symbol, true
	}

	// Look in outer environments
	if e.outer != nil {
		if symbol, ok := e.outer.Get(name); ok {
			return symbol, true
		}
	}

	// Fallback: check builtins metadata
	if info, ok := builtinTypes[name]; ok {
		ftype := &FunctionType{
			ParamTypes: info.ParamTypes,
			ReturnType: info.ReturnType,
		}
		return &Symbol{
			Name: name,
			Type: ftype,
		}, true
	}

	return nil, false
}

func (e *Environment) Set(name string, symbol *Symbol) {
	e.store[name] = symbol
}

// Type Checker
type TypeChecker struct {
	env           *Environment
	errors        []*TypeError
	currentReturn Type // The expected return type of the enclosing function
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

func (tc *TypeChecker) typeFromExpression(expr ast.Expression) Type {
	switch e := expr.(type) {
	case *ast.NumberLiteral:
		return &NumberType{}
	case *ast.StringLiteral:
		return &StringType{}
	case *ast.BooleanLiteral:
		return &BoolType{}
	case *ast.FunctionLiteral:
		paramTypes := []Type{}
		for _, param := range e.Parameters {
			paramTypes = append(paramTypes, tc.parseTypeFromIdentifier(param.TypeHint))
		}
		retType := tc.parseTypeFromIdentifier(e.ReturnType)
		return &FunctionType{
			ParamTypes: paramTypes,
			ReturnType: retType,
		}
	default:
		return &UnknownType{}
	}
}

// Helper function to parse type annotations from identifiers
func (tc *TypeChecker) parseTypeFromIdentifier(typeIdent *ast.Identifier) Type {
	if typeIdent == nil {
		return &UnknownType{}
	}

	switch typeIdent.Value {
	case NUMBER:
		return &NumberType{}
	case STRING:
		return &StringType{}
	case BOOLEAN:
		return &BoolType{}
	default:
		tc.addError(fmt.Sprintf("unknown type: %s", typeIdent.Value), typeIdent.Token.Line, typeIdent.Token.Column)
		return &UnknownType{}
	}
}

func GetTypeFromIdentifier(typeIdent *ast.Identifier) Type {
	if typeIdent == nil {
		return &UnknownType{}
	}

	switch typeIdent.Value {
	case NUMBER:
		return &NumberType{}
	case STRING:
		return &StringType{}
	case BOOLEAN:
		return &BoolType{}
	default:
		return &UnknownType{}
	}
}

// Main entry point for type checking
func (tc *TypeChecker) CheckProgram(program *ast.Program) Type {
	var last Type = &VoidType{}

	for _, stmt := range program.Statements {
		last = tc.CheckStatement(stmt)
	}
	return last
}

func (tc *TypeChecker) CheckStatement(stmt ast.Statement) Type {
	switch s := stmt.(type) {
	case *ast.LetStatement:
		return tc.CheckLetStatement(s)
	case *ast.ReturnStatement:
		return tc.CheckReturnStatement(s)
	case *ast.ExpressionStatement:
		return tc.CheckExpressionStatement(s)
	default:
		tc.addError(fmt.Sprintf("unknown statement type: %T", stmt), 0, 0)
		return &UnknownType{}
	}
}

func (tc *TypeChecker) CheckLetStatement(stmt *ast.LetStatement) Type {
	var declaredType Type
	if stmt.TypeHint != nil {
		declaredType = tc.parseTypeFromIdentifier(stmt.TypeHint)
	} else if stmt.Value != nil {
		declaredType = tc.typeFromExpression(stmt.Value)
	} else {
		tc.addError("type annotation is required for variable declarations", stmt.Token.Line, stmt.Token.Column)
		return &UnknownType{}
	}

	valueType := tc.CheckExpression(stmt.Value)

	if !declaredType.Equals(valueType) {
		tc.addError(
			fmt.Sprintf("type mismatch: declared %s but got %s", declaredType.String(), valueType.String()),
			stmt.Token.Line, stmt.Token.Column,
		)
	}

	tc.env.Set(stmt.Name.Value, &Symbol{
		Name:   stmt.Name.Value,
		Type:   declaredType,
		Line:   stmt.Name.Token.Line,
		Column: stmt.Name.Token.Column,
	})

	return &VoidType{}
}

func (tc *TypeChecker) CheckReturnStatement(stmt *ast.ReturnStatement) Type {
	if tc.currentReturn == nil {
		tc.addError("return statement outside of function", stmt.Token.Line, stmt.Token.Column)
		return &UnknownType{}
	}

	exprType := tc.CheckExpression(stmt.ReturnValue)

	if !tc.currentReturn.Equals(exprType) {
		tc.addError(
			fmt.Sprintf("return type mismatch: expected %s, got %s", tc.currentReturn.String(), exprType.String()),
			stmt.Token.Line,
			stmt.Token.Column,
		)
	}

	return exprType // <- Return the type of the returned expression
}

func (tc *TypeChecker) CheckExpressionStatement(stmt *ast.ExpressionStatement) Type {
	if stmt.Expression == nil {
		tc.addError("empty expression statement", stmt.Token.Line, stmt.Token.Column)
		return &UnknownType{}
	}

	exprType := tc.CheckExpression(stmt.Expression)

	if stmt.HasSemicolon {
		return &VoidType{} // semicolon → discard value
	}

	return exprType // no semicolon → use expression type
}

func (tc *TypeChecker) CheckExpression(expr ast.Expression) Type {
	switch e := expr.(type) {
	case *ast.NumberLiteral:
		return &NumberType{}
	case *ast.StringLiteral:
		return &StringType{}
	case *ast.BooleanLiteral:
		return &BoolType{}
	case *ast.Identifier:
		return tc.CheckIdentifier(e)
	case *ast.InfixExpression:
		return tc.CheckInfixExpression(e)
	case *ast.PrefixExpression:
		return tc.CheckPrefixExpression(e)
	case *ast.IfExpression:
		return tc.CheckIfExpression(e)
	case *ast.FunctionLiteral:
		return tc.CheckFunctionLiteral(e)
	case *ast.CallExpression:
		return tc.CheckCallExpression(e)
	case *ast.AssignmentExpression:
		return tc.CheckAssignmentExpression(e)
	default:
		tc.addError(fmt.Sprintf("unknown expression type: %T", expr), 0, 0)
		return &UnknownType{}
	}
}

func (tc *TypeChecker) CheckIdentifier(ident *ast.Identifier) Type {
	symbol, exists := tc.env.Get(ident.Value)
	if !exists {
		tc.addError(fmt.Sprintf("undefined variable: %s", ident.Value), ident.Token.Line, ident.Token.Column)
		return &UnknownType{}
	}
	return symbol.Type
}

func (tc *TypeChecker) CheckInfixExpression(expr *ast.InfixExpression) Type {
	leftType := tc.CheckExpression(expr.Left)
	rightType := tc.CheckExpression(expr.Right)

	// If either side is UnknownType, propagate UnknownType instead of adding errors
	if _, ok := leftType.(*UnknownType); ok {
		return &UnknownType{}
	}
	if _, ok := rightType.(*UnknownType); ok {
		return &UnknownType{}
	}

	switch expr.Operator {
	case "+":
		// If both sides are numbers, result is Number
		if leftType.Equals(&NumberType{}) && rightType.Equals(&NumberType{}) {
			return &NumberType{}
		}

		// If both sides are strings, result is String
		if leftType.Equals(&StringType{}) && rightType.Equals(&StringType{}) {
			return &StringType{}
		}

		// Anything else is an error
		tc.addError(fmt.Sprintf("cannot add %s and %s", leftType.String(), rightType.String()), expr.Token.Line, expr.Token.Column)
		return &UnknownType{}

	case "-", "*", "/":
		// Arithmetic operators require numbers
		if !leftType.Equals(&NumberType{}) {
			tc.addError(
				fmt.Sprintf(
					"left operand of %s must be %s, got %s",
					expr.Operator,
					NUMBER,
					leftType.String()),
				expr.Token.Line,
				expr.Token.Column)
		}

		if !rightType.Equals(&NumberType{}) {
			tc.addError(
				fmt.Sprintf(
					"right operand of %s must be %s, got %s",
					expr.Operator,
					NUMBER,
					rightType.String()),
				expr.Token.Line,
				expr.Token.Column)
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
			tc.addError(
				fmt.Sprintf(
					"left operand of %s must be %s, got %s",
					expr.Operator,
					NUMBER,
					leftType.String()),
				expr.Token.Line,
				expr.Token.Column)
		}

		if !rightType.Equals(&NumberType{}) {
			tc.addError(
				fmt.Sprintf(
					"right operand of %s must be %s, got %s",
					expr.Operator,
					NUMBER,
					rightType.String()),
				expr.Token.Line,
				expr.Token.Column)
		}
		return &BoolType{}

	default:
		tc.addError(fmt.Sprintf("unknown infix operator: %s", expr.Operator), expr.Token.Line, expr.Token.Column)
		return &UnknownType{}
	}
}

func (tc *TypeChecker) CheckPrefixExpression(expr *ast.PrefixExpression) Type {
	operandType := tc.CheckExpression(expr.Right)

	// Propagate up instead of adding more errors.
	if _, ok := operandType.(*UnknownType); ok {
		return &UnknownType{}
	}

	switch expr.Operator {
	case "-":
		if !operandType.Equals(&NumberType{}) {
			tc.addError(fmt.Sprintf("unary minus requires %s, got %s", NUMBER, operandType.String()), expr.Token.Line, expr.Token.Column)
		}
		return &NumberType{}

	case "!":
		if !operandType.Equals(&BoolType{}) {
			tc.addError(fmt.Sprintf("logical not requires %s, got %s", BOOLEAN, operandType.String()), expr.Token.Line, expr.Token.Column)
		}
		return &BoolType{}

	default:
		tc.addError(fmt.Sprintf("unknown prefix operator: %s", expr.Operator), expr.Token.Line, expr.Token.Column)
		return &UnknownType{}
	}
}

func (tc *TypeChecker) CheckIfExpression(expr *ast.IfExpression) Type {
	condType := tc.CheckExpression(expr.Condition)

	// Condition must be Boolean
	if !condType.Equals(&BoolType{}) {
		tc.addError(fmt.Sprintf("if condition must be %s, got %s", BOOLEAN, condType.String()), expr.Token.Line, expr.Token.Column)
	}

	// Check consequence block
	consequenceType := tc.CheckBlockStatement(expr.Consequence)

	// Check alternative block if present
	var alternativeType Type = &VoidType{}
	if expr.Alternative != nil {
		alternativeType = tc.CheckBlockStatement(expr.Alternative)
	}

	// Both branches must have the same type if alternative exists
	if expr.Alternative != nil && !consequenceType.Equals(alternativeType) {
		tc.addError(fmt.Sprintf(
			"if branches must return same type, got %s and %s",
			consequenceType.String(),
			alternativeType.String(),
		), expr.Token.Line, expr.Token.Column)
		return &UnknownType{}
	}

	if expr.Alternative != nil {
		return consequenceType
	}
	return &VoidType{} // If no else branch, type is Void
}

func (tc *TypeChecker) CheckAssignmentExpression(expr *ast.AssignmentExpression) Type {
	sym, ok := tc.env.Get(expr.Name.Value)
	if !ok {
		tc.addError(fmt.Sprintf("variable with name %s not defined", expr.Name.Value), expr.Token.Line, expr.Token.Column)
		return &UnknownType{}
	}
	newType := tc.CheckExpression(expr.Value)

	if !sym.Type.Equals(newType) {
		tc.addError(fmt.Sprintf("assignment type mismatch, expected %s, got %s", sym.Type, newType), expr.Token.Line, expr.Token.Line)
		return &UnknownType{}
	}

	return sym.Type
}

func (tc *TypeChecker) CheckBlockStatement(block *ast.BlockStatement) Type {
	var lastType Type = &VoidType{}
	for _, stmt := range block.Statements {
		lastType = tc.CheckStatement(stmt)
	}
	return lastType
}

func (tc *TypeChecker) CheckFunctionLiteral(fn *ast.FunctionLiteral) Type {
	// Parse parameter types
	paramTypes := []Type{}
	for _, param := range fn.Parameters {
		paramTypes = append(paramTypes, tc.parseTypeFromIdentifier(param.TypeHint))
	}

	// Expected return type
	var returnType Type
	returnType = &UnknownType{}
	if fn.ReturnType != nil {
		returnType = tc.parseTypeFromIdentifier(fn.ReturnType)
	}

	// --- Predeclare function in current environment ---
	fnType := &FunctionType{
		ParamTypes: paramTypes,
		ReturnType: returnType,
	}

	// If the function has a name, insert it in the environment now
	tc.env.Set(fn.Name, &Symbol{
		Name:   fn.Name,
		Type:   fnType,
		Line:   fn.Token.Line,
		Column: fn.Token.Column,
	})

	// New scope for function body
	oldEnv := tc.env
	oldReturn := tc.currentReturn
	tc.env = NewEnclosedEnvironment(oldEnv)
	tc.currentReturn = returnType

	// Add parameters to environment
	for i, param := range fn.Parameters {
		tc.env.Set(param.Name.Value, &Symbol{
			Name:   param.Name.Value,
			Type:   paramTypes[i],
			Line:   param.Name.Token.Line,
			Column: param.Name.Token.Column,
		})
	}

	// Check body
	bodyType := tc.CheckBlockStatement(fn.Body)

	// Restore old environment
	tc.env = oldEnv
	tc.currentReturn = oldReturn

	// Ensure body type matches declared return type
	if !returnType.Equals(bodyType) {
		tc.addError(fmt.Sprintf(
			"function body type mismatch: expected %s, got %s",
			returnType.String(), bodyType.String(),
		), fn.Token.Line, fn.Token.Column)
	}

	// Return a FunctionType (you’ll need to define it)
	return &FunctionType{
		ParamTypes: paramTypes,
		ReturnType: returnType,
	}
}

func (tc *TypeChecker) CheckCallExpression(ce *ast.CallExpression) Type {
	fnType := tc.CheckExpression(ce.Function)

	if ident, ok := ce.Function.(*ast.Identifier); ok {
		if info, exists := builtinTypes[ident.Value]; exists {
			// check arity
			if info.Arity != -1 && len(ce.Arguments) != info.Arity {
				tc.addError(fmt.Sprintf("argument count mismatch: expected %d, got %d", info.Arity, len(ce.Arguments)), ce.Token.Line, ce.Token.Column)
				return &UnknownType{}
			}

			// check argument types
			// If arity is -1, then we only expect a single param type and all values must be that param type.
			for i, arg := range ce.Arguments {
				argType := tc.CheckExpression(arg)
				if info.Arity == -1 && !info.ParamTypes[0].Equals(&UnknownType{}) && !argType.Equals(info.ParamTypes[0]) {
					tc.addError(fmt.Sprintf("argument %d type mismatch: expected %v, got %v", i+1, info.ParamTypes[0], argType), ce.Token.Line, ce.Token.Column)
					return &UnknownType{}
				}

				if info.Arity != -1 && !info.ParamTypes[i].Equals(&UnknownType{}) && !argType.Equals(info.ParamTypes[i]) {
					tc.addError(fmt.Sprintf("argument %d type mismatch: expected %v, got %v", i+1, info.ParamTypes[i], argType), ce.Token.Line, ce.Token.Column)
					return &UnknownType{}
				}
			}

			return info.ReturnType
		}
	}

	fn, ok := fnType.(*FunctionType)
	if !ok {
		tc.addError(fmt.Sprintf("attempted to call a non-function type: %v", fnType), ce.Token.Line, ce.Token.Column)
		return &UnknownType{}
	}

	if len(ce.Arguments) != len(fn.ParamTypes) {
		tc.addError(fmt.Sprintf("argument count mismatch: expected %d, got %d", len(fn.ParamTypes), len(ce.Arguments)), ce.Token.Line, ce.Token.Column)
		return &UnknownType{}
	}

	for i, arg := range ce.Arguments {
		argType := tc.CheckExpression(arg)
		if !argType.Equals(fn.ParamTypes[i]) {
			tc.addError(fmt.Sprintf("argument %d type mismatch: expected %v, got %v", i+1, fn.ParamTypes[i], argType), ce.Token.Line, ce.Token.Column)
			return &UnknownType{}
		}
	}

	return fn.ReturnType
}
