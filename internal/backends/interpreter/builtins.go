package interpreter

import (
	"fmt"
	"sigil/internal/typechecker"
)

var builtins = map[string]*Builtin{
	"len": {
		Arity: 1,
		Fn: func(args ...Value) (Value, error) {
			switch a := args[0].(type) {
			case *StringValue:
				return &NumberValue{Value: float64(len(a.Value))}, nil
			default:
				return nil, fmt.Errorf("len not defined for type %s", a.Type())
			}
		},
		ParamTypes: []typechecker.Type{&typechecker.StringType{}},
		ReturnType: &typechecker.NumberType{},
	},
	"print": {
		Arity: -1, // variadic
		Fn: func(args ...Value) (Value, error) {
			for i, arg := range args {
				s, ok := arg.(*StringValue)
				if !ok {
					return nil, fmt.Errorf("print only accepts strings, got %s", arg.Type())
				}
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(s.Value)
			}
			return &VoidValue{}, nil
		},
		ReturnType: &typechecker.VoidType{},
		ParamTypes: nil, // variadic, only String allowed at runtime
	},
	"println": {
		Arity: -1, // variadic
		Fn: func(args ...Value) (Value, error) {
			for i, arg := range args {
				s, ok := arg.(*StringValue)
				if !ok {
					return nil, fmt.Errorf("println only accepts strings, got %s", arg.Type())
				}
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(s.Value)
			}
			fmt.Println()
			return &VoidValue{}, nil
		},
		ReturnType: &typechecker.VoidType{},
		ParamTypes: nil,
	},
}
