package interpreter

import (
	"fmt"
)

var builtins = map[string]*Builtin{
	"len": {
		Name:  "len",
		Arity: 1,
		Fn: func(args ...Value) (Value, error) {
			switch a := args[0].(type) {
			case *StringValue:
				return &NumberValue{Value: float64(len(a.Value))}, nil
			default:
				return nil, fmt.Errorf("len not defined for type %s", a.Type())
			}
		},
	},
	"print": {
		Name:  "print",
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
	},
	"println": {
		Name:  "println",
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
	},
	"string": {
		Name:  "string",
		Arity: 1,
		Fn: func(args ...Value) (Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("expected 1 argument, but got %d", len(args))
			}

			return &StringValue{Value: args[0].String()}, nil
		},
	},
}
