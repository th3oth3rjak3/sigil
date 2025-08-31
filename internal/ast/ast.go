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

// LetStatement represents variable declarations
type LetStatement struct {
	Token    lexer.Token // the LET token
	Name     *Identifier
	TypeHint *Identifier // for type annotations like ': Number'
	Value    Expression
}

func (ls *LetStatement) stmt()                {}
func (ls *LetStatement) TokenLiteral() string { return "let" }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	if ls.TypeHint != nil {
		out.WriteString(": " + ls.TypeHint.String())
	}
	out.WriteString(" = " + ls.Value.String() + ";")
	return out.String()
}

func (ls *LetStatement) TreeString(prefix string, isLast bool) string {
	var out strings.Builder

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	out.WriteString(prefix + connector + "LetStatement\n")

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	out.WriteString(childPrefix + "├── Name: " + ls.Name.String() + "\n")

	if ls.TypeHint != nil {
		out.WriteString(childPrefix + "├── TypeHint: " + ls.TypeHint.String() + "\n")
	}

	out.WriteString(childPrefix + "└── Value: " + ls.Value.String() + "\n")

	return out.String()
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
	Value string
}

func (nl *NumberLiteral) expr()                {}
func (nl *NumberLiteral) String() string       { return nl.Value }
func (nl *NumberLiteral) TokenLiteral() string { return nl.Value }

func (nl *NumberLiteral) TreeString(prefix string, isLast bool) string {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	return prefix + connector + "NumberLiteral: " + nl.Value + "\n"
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

// InfixExpression represents binary operations like a + b
type InfixExpression struct {
	Token    lexer.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expr()                {}
func (ie *InfixExpression) TokenLiteral() string { return "TODO" }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

func (ie *InfixExpression) TreeString(prefix string, isLast bool) string {
	var out strings.Builder

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	out.WriteString(prefix + connector + "InfixExpression: " + ie.Operator + "\n")

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	out.WriteString(ie.Left.TreeString(childPrefix, false))
	out.WriteString(ie.Right.TreeString(childPrefix, true))

	return out.String()
}

// PrefixExpression represents unary operations like -a or !b
type PrefixExpression struct {
	Token    lexer.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expr()                {}
func (pe *PrefixExpression) TokenLiteral() string { return "TODO" }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

func (pe *PrefixExpression) TreeString(prefix string, isLast bool) string {
	var out strings.Builder

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	out.WriteString(prefix + connector + "PrefixExpression: " + pe.Operator + "\n")

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	out.WriteString(pe.Right.TreeString(childPrefix, true))

	return out.String()
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

// ReturnStatement is a return value from a function.
type ReturnStatement struct {
	Token       lexer.Token // The "return" token
	ReturnValue Expression  // The expression value to return
}

func (rs *ReturnStatement) stmt()                {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
func (rs *ReturnStatement) TreeString(prefix string, isLast bool) string {
	var out strings.Builder

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	out.WriteString(prefix + connector + "ReturnStatement: " + rs.TokenLiteral() + "\n")

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.TreeString(childPrefix, true))
	}

	return out.String()
}

type ExpressionStatement struct {
	Token      lexer.Token // The first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) stmt()                {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}
func (es *ExpressionStatement) TreeString(prefix string, isLast bool) string {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	return prefix + connector + "ExpressionStatement: " + es.String() + "\n"
}
