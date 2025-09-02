package typechecker

import (
	"fmt"
	"strings"
)

const (
	NUMBER  = "Number"
	BOOLEAN = "Boolean"
	STRING  = "String"
	VOID    = "Void"
	UNKNOWN = "Unknown"
)

// Type represents a type in the language
type Type interface {
	String() string
	Equals(Type) bool
}

type BuiltinTypeInfo struct {
	Arity      int // when -1, then variadic
	ParamTypes []Type
	ReturnType Type
}

// Basic types
type NumberType struct{}
type StringType struct{}
type BoolType struct{}
type VoidType struct{}    // for statements that don't return values
type UnknownType struct{} // For errors during type checking

func (nt *NumberType) String() string { return NUMBER }
func (nt *NumberType) Equals(other Type) bool {
	_, ok := other.(*NumberType)
	return ok
}

func (st *StringType) String() string { return STRING }
func (st *StringType) Equals(other Type) bool {
	_, ok := other.(*StringType)
	return ok
}

func (bt *BoolType) String() string { return BOOLEAN }
func (bt *BoolType) Equals(other Type) bool {
	_, ok := other.(*BoolType)
	return ok
}

func (vt *VoidType) String() string { return VOID }
func (vt *VoidType) Equals(other Type) bool {
	_, ok := other.(*VoidType)
	return ok
}

func (ut *UnknownType) String() string { return UNKNOWN }
func (ut *UnknownType) Equals(other Type) bool {
	_, ok := other.(*UnknownType)
	return ok
}

// FunctionType represents a function with parameters and a return type
type FunctionType struct {
	ParamTypes []Type
	ReturnType Type
}

func (ft *FunctionType) String() string {
	params := []string{}
	for _, p := range ft.ParamTypes {
		params = append(params, p.String())
	}
	return fmt.Sprintf("(%s) -> %s", strings.Join(params, ", "), ft.ReturnType.String())
}

func (ft *FunctionType) Equals(other Type) bool {
	o, ok := other.(*FunctionType)
	if !ok {
		return false
	}
	if len(ft.ParamTypes) != len(o.ParamTypes) {
		return false
	}
	for i := range ft.ParamTypes {
		if !ft.ParamTypes[i].Equals(o.ParamTypes[i]) {
			return false
		}
	}
	return ft.ReturnType.Equals(o.ReturnType)
}

// Type error with position information
type TypeError struct {
	Message string
	Line    int
	Column  int
}

func (te *TypeError) Error() string {
	return fmt.Sprintf("Type error at line %d, column %d: %s", te.Line, te.Column, te.Message)
}
