package interpreter

import (
	"math"
)

const EPSILON = 1e-9

func evalInfixExpression(operator string, left, right Object) Object {
	if left.Type() != right.Type() {
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}

	switch {
	case left.Type() == NUMBER_OBJ && right.Type() == NUMBER_OBJ:
		return evalNumberInfixExpression(operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ && operator == "+":
		return evalConcatStrings(left, right)

	// IMPORTANT NOTE: these do pointer comparison because we assume
	// that if we've reached this point, any value comparisons have come
	// before this in the switch statement. Put new value comparisons
	// up there ^
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalNumberInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Number).Value
	rightVal := right.(*Number).Value

	switch operator {
	// Arithmetic
	case "+":
		return &Number{Value: leftVal + rightVal}
	case "-":
		return &Number{Value: leftVal - rightVal}
	case "*":
		return &Number{Value: leftVal * rightVal}
	case "/":
		return &Number{Value: leftVal / rightVal}

	// Comparison
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)

	// Equality
	case "==":
		return nativeBoolToBooleanObject(math.Abs(leftVal-rightVal) <= EPSILON)
	case "!=":
		return nativeBoolToBooleanObject(math.Abs(leftVal-rightVal) > EPSILON)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalConcatStrings(left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value

	return &String{Value: leftVal + rightVal}
}
