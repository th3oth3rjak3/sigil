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
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	var out strings.Builder
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
	// Value (nested tree)
	out.WriteString(childPrefix + "└── Value:\n")
	out.WriteString(ls.Value.TreeString(childPrefix+"    ", true))

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
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	var out strings.Builder
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
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	var out strings.Builder
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

	out.WriteString(prefix + connector + "ReturnStatement\n")

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	if rs.ReturnValue != nil {
		// Recursively print the value as a child
		out.WriteString(rs.ReturnValue.TreeString(childPrefix, true))
	}

	return out.String()
}

type ExpressionStatement struct {
	Token        lexer.Token // The first token of the expression
	Expression   Expression
	HasSemicolon bool
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
	var out strings.Builder
	out.WriteString(prefix + connector + "ExpressionStatement\n")

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	if es.Expression != nil {
		out.WriteString(es.Expression.TreeString(childPrefix, true))
	}

	return out.String()
}

type BlockStatement struct {
	Token      lexer.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) stmt() {}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) TreeString(prefix string, isLast bool) string {
	var out strings.Builder

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	out.WriteString(prefix + connector + "BlockStatement\n")

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	for i, stmt := range bs.Statements {
		stmtIsLast := i == len(bs.Statements)-1
		out.WriteString(stmt.TreeString(childPrefix, stmtIsLast))
	}

	return out.String()
}

type IfExpression struct {
	Token       lexer.Token     // The 'if' token
	Condition   Expression      // The condition to evaluate to decide which branch to take
	Consequence *BlockStatement // Taken if condition is true
	Alternative *BlockStatement // Taken if condition is false
}

func (ie *IfExpression) expr() {}
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) TreeString(prefix string, isLast bool) string {
	var out strings.Builder

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	out.WriteString(prefix + connector + "IfExpression\n")

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Condition
	out.WriteString(childPrefix + "├── Condition:\n")
	out.WriteString(ie.Condition.TreeString(childPrefix+"│   ", true))

	// Consequence
	out.WriteString(childPrefix + "├── Consequence:\n")
	out.WriteString(ie.Consequence.TreeString(childPrefix+"│   ", true))

	// Alternative (optional)
	if ie.Alternative != nil {
		out.WriteString(childPrefix + "└── Alternative:\n")
		out.WriteString(ie.Alternative.TreeString(childPrefix+"    ", true))
	}

	return out.String()
}

type FunctionParameter struct {
	Name     *Identifier
	TypeHint *Identifier // may be nil if no type hint provided
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
	ReturnType *Identifier // optional return type
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

type CallExpression struct {
	Token     lexer.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expr()                {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (ce *CallExpression) TreeString(prefix string, isLast bool) string {
	var out strings.Builder

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	out.WriteString(prefix + connector + "CallExpression\n")

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Function
	out.WriteString(childPrefix + "├── Function:\n")
	out.WriteString(ce.Function.TreeString(childPrefix+"│   ", true))

	// Arguments
	out.WriteString(childPrefix + "└── Arguments:\n")
	for i, arg := range ce.Arguments {
		argIsLast := i == len(ce.Arguments)-1
		out.WriteString(arg.TreeString(childPrefix+"    ", argIsLast))
	}

	return out.String()
}
