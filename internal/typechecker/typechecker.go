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

func (tc *TypeChecker) parseTypeFromAstType(t ast.Type) Type {
	switch tt := t.(type) {
	case *ast.SimpleType:
		switch tt.Name {
		case NUMBER:
			return &NumberType{}
		case STRING:
			return &StringType{}
		case BOOLEAN:
			return &BoolType{}
		case VOID:
			return &VoidType{}
		default:
			tc.addError(fmt.Sprintf("unknown type: %s", tt.Name), tt.Token.Line, tt.Token.Column)
			return &UnknownType{}
		}
	case *ast.FunctionType:
		paramTypes := []Type{}
		for _, p := range tt.ParamTypes {
			paramTypes = append(paramTypes, tc.parseTypeFromAstType(p))
		}
		retType := tc.parseTypeFromAstType(tt.ReturnType)
		return &FunctionType{
			ParamTypes: paramTypes,
			ReturnType: retType,
		}
	default:
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
