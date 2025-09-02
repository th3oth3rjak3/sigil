package interpreter

import (
	"fmt"
	"strconv"
)

var (
	NULL  Object = &Null{}
	TRUE  Object = &Boolean{Value: true}
	FALSE Object = &Boolean{Value: false}
)

const (
	NUMBER_OBJ  = "Number"
	BOOLEAN_OBJ = "Boolean"
	STRING_OBJ  = "String"
	NULL_OBJ    = "Null"
	RETURN_OBJ  = "RETURN_OBJ"
	ERROR_OBJ   = "Error"
)

type ObjectType string

// Object is the interpreter base type.
type Object interface {
	// Get the type of the object.
	Type() ObjectType
	// Inspect the object contents.
	Inspect() string
}

// Number represents all real numbers including
// floating point and integer numbers
type Number struct {
	// The actual value.
	Value float64
}

func (n *Number) Inspect() string  { return strconv.FormatFloat(n.Value, 'f', -1, 64) }
func (n *Number) Type() ObjectType { return NUMBER_OBJ }

// Boolean represents the values true and false.
type Boolean struct {
	// The actual value.
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// The absence of a meaningful value.
type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }

type ReturnObject struct {
	Value Object
}

func (ro *ReturnObject) Inspect() string  { return ro.Value.Inspect() }
func (ro *ReturnObject) Type() ObjectType { return RETURN_OBJ }

type Error struct {
	Message string
}

func (e *Error) Inspect() string  { return "Error: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }

func newError(format string, args ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}

func isError(obj Object) bool {
	return obj != nil && obj.Type() == ERROR_OBJ
}
