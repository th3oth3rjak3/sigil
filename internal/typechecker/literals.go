package typechecker

import (
	"fmt"
	"sigil/internal/ast"
)

func (tc *TypeChecker) CheckFunctionLiteral(fn *ast.FunctionLiteral) Type {
	// Parse parameter types
	paramTypes := []Type{}
	for _, param := range fn.Parameters {
		paramTypes = append(paramTypes, tc.parseTypeFromAstType(param.TypeHint))
	}

	// Expected return type
	var returnType Type
	returnType = &UnknownType{}
	if fn.ReturnType != nil {
		returnType = tc.parseTypeFromAstType(fn.ReturnType)
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
