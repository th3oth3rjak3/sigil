package interpreter

import "sigil/internal/ast"

func nativeBoolToBooleanObject(input bool) Object {
	if input {
		return TRUE
	}

	return FALSE
}

func evalIfExpression(expr *ast.IfExpression) Object {
	condition := Eval(expr.Condition)
	if isError(condition) {
		return condition
	}

	// We only allow booleans as conditionals.
	// This should be validated in the typechecker well before
	// we get here.
	condVal, ok := condition.(*Boolean)
	if !ok {
		return newError("type mismatch: expected %s but got %s", BOOLEAN_OBJ, condition.Type())
	}

	if condVal.Value {
		return Eval(expr.Consequence)
	} else if expr.Alternative != nil {
		return Eval(expr.Alternative)
	} else {
		return NULL
	}
}
