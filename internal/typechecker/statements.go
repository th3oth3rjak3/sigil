package typechecker

import (
	"fmt"
	"sigil/internal/ast"
)

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

func (tc *TypeChecker) CheckBlockStatement(block *ast.BlockStatement) Type {
	var lastType Type = &VoidType{}
	for _, stmt := range block.Statements {
		lastType = tc.CheckStatement(stmt)
	}
	return lastType
}
