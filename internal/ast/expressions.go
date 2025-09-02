package ast

import (
	"sigil/internal/lexer"
	"strings"
)

type AssignmentExpression struct {
	Token lexer.Token // the '=' token
	Name  *Identifier
	Value Expression
}

func (ae *AssignmentExpression) expr()                {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	var out strings.Builder
	out.WriteString(ae.Name.Value)
	out.WriteString(" = ")
	out.WriteString(ae.Value.TokenLiteral())
	out.WriteString(";")

	return out.String()
}
func (ae *AssignmentExpression) TreeString(prefix string, isLast bool) string {
	var out strings.Builder

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	// Root line for this assignment
	out.WriteString(prefix)
	out.WriteString(connector)
	out.WriteString("AssignmentExpression\n")

	// New prefix for children
	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Name
	out.WriteString(childPrefix)
	out.WriteString("├── Name: ")
	out.WriteString(ae.Name.Value)
	out.WriteString("\n")

	// Value
	out.WriteString(childPrefix)
	out.WriteString("└── Value:\n")
	out.WriteString(ae.Value.TreeString(childPrefix+"    ", true))

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
	var out strings.Builder

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

type IfExpression struct {
	Token       lexer.Token     // The 'if' token
	Condition   Expression      // The condition to evaluate to decide which branch to take
	Consequence *BlockStatement // Taken if condition is true
	Alternative *BlockStatement // Taken if condition is false
}

func (ie *IfExpression) expr() {}
func (ie *IfExpression) String() string {
	var out strings.Builder

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
