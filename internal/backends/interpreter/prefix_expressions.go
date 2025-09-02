package interpreter

// Prefix Expressions
func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right Object) Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return NULL
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != NUMBER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*Number).Value
	return &Number{Value: -value}
}
