package interpreter

import "sigil/internal/ast"

func evalProgram(program *ast.Program, env *EvaluatorEnvironment) Object {
	var result Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *ReturnObject:
			return result.Value
		case *Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *EvaluatorEnvironment) Object {
	var result Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == RETURN_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}

	return result
}
