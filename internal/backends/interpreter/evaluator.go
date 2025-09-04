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
	env := NewEvaluatorEnvironment()
	for _, stmt := range program.Statements {
		val := Eval(stmt, env)

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
func Eval(node ast.Node, env *EvaluatorEnvironment) Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		// Not returning the value because
		// we don't allow things like this: y = x = 10
		// Returning this value might encourage bad behavior.
		_ = env.Set(node.Name.Value, val)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
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

	case *ast.FunctionLiteral:
		return &Function{
			Name:       node.Name,
			Parameters: node.Parameters,
			Body:       node.Body,
			ReturnType: node.ReturnType,
			Env:        env,
		}

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.AssignmentExpression:
		right := Eval(node.Value, env)
		if isError(right) {
			return right
		}

		ok := false
		e := env
		for !ok {
			ok = e.Contains(node.Name.Value)
			if ok {
				e.Set(node.Name.Value, right)
				break
			} else if e.outer != nil {
				e = e.outer
			} else {
				return newError("undefined variable: %s", node.Name.Value)
			}
		}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	}

	return nil
}

func applyFunction(fun Object, args []Object) Object {
	function, ok := fun.(*Function)
	if !ok {
		return newError("not a function: %s", fun.Type())
	}

	extendedEnv := extendFunctionEnvironment(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnvironment(fun *Function, args []Object) *EvaluatorEnvironment {
	env := NewEnclosedEvaluatorEnvironment(fun.Env)

	for paramIdx, param := range fun.Parameters {
		env.Set(param.Name.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnObject); ok {
		return returnValue.Value
	}

	return obj
}
