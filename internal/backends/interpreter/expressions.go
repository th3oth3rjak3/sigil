package interpreter

import "sigil/internal/ast"

func nativeBoolToBooleanObject(input bool) Object {
	if input {
		return TRUE
	}

	return FALSE
}

func evalIfExpression(expr *ast.IfExpression, env *EvaluatorEnvironment) Object {
	condition := Eval(expr.Condition, env)
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
		return Eval(expr.Consequence, env)
	} else if expr.Alternative != nil {
		return Eval(expr.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(expr *ast.Identifier, env *EvaluatorEnvironment) Object {
	val, ok := env.Get(expr.Value)
	if !ok {
		return newError("identifier not found: %s", expr.Value)
	}

	return val
}

func evalExpressions(exps []ast.Expression, env *EvaluatorEnvironment) []Object {
	var result []Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}
