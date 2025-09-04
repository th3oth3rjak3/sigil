package ast

import (
	"bytes"
	"fmt"
	"sigil/internal/lexer"
	"strings"
)

// Node represents any node in the AST
type Node interface {
	TokenLiteral() string
	String() string
	TreeString(prefix string, isLast bool) string
}

// Statement represents statement nodes
type Statement interface {
	Node
	stmt()
}

// Expression represents expression nodes
type Expression interface {
	Node
	expr()
}

// Program is the root node of every AST
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) TreeString(prefix string, isLast bool) string {
	var out strings.Builder
	out.WriteString("Program\n")
	for i, stmt := range p.Statements {
		isLast := i == len(p.Statements)-1
		out.WriteString(stmt.TreeString("", isLast))
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// Identifier represents variable names, type names, etc.
type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expr()                {}
func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) TokenLiteral() string { return i.Value }

func (i *Identifier) TreeString(prefix string, isLast bool) string {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	return prefix + connector + "Identifier: " + i.Value + "\n"
}

// NumberLiteral represents numeric values
type NumberLiteral struct {
	Token lexer.Token
	Value float64
}

func (nl *NumberLiteral) expr()                {}
func (nl *NumberLiteral) String() string       { return nl.Token.Literal }
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Literal }

func (nl *NumberLiteral) TreeString(prefix string, isLast bool) string {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	return prefix + connector + "NumberLiteral: " + nl.TokenLiteral() + "\n"
}

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expr()                {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Value }
func (sl *StringLiteral) String() string       { return sl.Value }

func (sl *StringLiteral) TreeString(prefix string, isLast bool) string {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	return prefix + connector + "StringLiteral: " + sl.Value + "\n"
}

// BooleanLiteral for true/false values
type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (bl *BooleanLiteral) expr()                {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.String() }
func (bl *BooleanLiteral) String() string {
	return fmt.Sprintf("%t", bl.Value)
}

func (bl *BooleanLiteral) TreeString(prefix string, isLast bool) string {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	return prefix + connector + "BooleanLiteral: " + bl.String() + "\n"
}

type FunctionParameter struct {
	Name     *Identifier
	TypeHint Type // may be nil if no type hint provided
}

func (fp *FunctionParameter) String() string {
	if fp.TypeHint != nil {
		return fp.Name.String() + ": " + fp.TypeHint.String()
	}
	return fp.Name.String()
}

func (fp *FunctionParameter) TreeString(prefix string, isLast bool) string {
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	var out strings.Builder
	out.WriteString(prefix + connector + "FunctionParameter\n")

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Name
	out.WriteString(childPrefix + "├── Name: " + fp.Name.String() + "\n")

	// TypeHint
	if fp.TypeHint != nil {
		out.WriteString(childPrefix + "└── TypeHint: " + fp.TypeHint.String() + "\n")
	} else {
		out.WriteString(childPrefix + "└── TypeHint: (none)\n")
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      lexer.Token // The 'fn' token
	Name       string
	Parameters []*FunctionParameter // a list of parameters, may be empty
	Body       *BlockStatement
	ReturnType Type
}

func (fl *FunctionLiteral) expr() {}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	// Function keyword
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")

	// Parameters
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	// Return type (if any)
	if fl.ReturnType != nil {
		out.WriteString(": ")
		out.WriteString(fl.ReturnType.String())
	}

	// Body
	out.WriteString(" ")
	out.WriteString(fl.Body.String())

	return out.String()
}

func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

func (fl *FunctionLiteral) TreeString(prefix string, isLast bool) string {
	var out strings.Builder

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	out.WriteString(prefix + connector + "FunctionLiteral\n")

	// Compute prefix for children
	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Parameters
	for i, param := range fl.Parameters {
		paramIsLast := i == len(fl.Parameters)-1 && fl.ReturnType == nil && fl.Body == nil
		out.WriteString(param.TreeString(childPrefix, paramIsLast))
	}

	// Return type
	if fl.ReturnType != nil {
		retIsLast := fl.Body == nil
		retConnector := "├── "
		if retIsLast {
			retConnector = "└── "
		}
		out.WriteString(childPrefix + retConnector + "ReturnType: " + fl.ReturnType.String() + "\n")
	}

	// Body
	if fl.Body != nil {
		bodyConnector := "└── "
		out.WriteString(childPrefix + bodyConnector + "Body:\n")
		out.WriteString(fl.Body.TreeString(childPrefix+"    ", true))
	}

	return out.String()
}

type Type interface {
	Node
	typeNode()
	expr()
	String() string
}

type SimpleType struct {
	Token lexer.Token
	Name  string
}

func (st *SimpleType) expr()                {}
func (st *SimpleType) typeNode()            {}
func (st *SimpleType) TokenLiteral() string { return st.Token.Literal }
func (st *SimpleType) String() string       { return st.Name }
func (st *SimpleType) TreeString(prefix string, isLast bool) string {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	return prefix + connector + "SimpleType: " + st.Name + "\n"
}

type FunctionType struct {
	ParamTypes []Type
	ReturnType Type
}

func (ft *FunctionType) expr()                {}
func (ft *FunctionType) typeNode()            {}
func (ft *FunctionType) TokenLiteral() string { return "fn" }
func (ft *FunctionType) String() string {
	parts := []string{}
	for _, p := range ft.ParamTypes {
		parts = append(parts, p.String())
	}
	return "(" + strings.Join(parts, ", ") + ") -> " + ft.ReturnType.String()
}
func (ft *FunctionType) TreeString(prefix string, isLast bool) string {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	var out strings.Builder
	out.WriteString(prefix + connector + "FunctionType\n")

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Parameter types
	for i, p := range ft.ParamTypes {
		paramIsLast := i == len(ft.ParamTypes)-1 && ft.ReturnType == nil
		out.WriteString(p.TreeString(childPrefix, paramIsLast))
	}

	// Return type
	if ft.ReturnType != nil {
		out.WriteString(ft.ReturnType.TreeString(childPrefix, true))
	}

	return out.String()
}
