package ast

import (
	"sigil/internal/lexer"
	"strings"
)

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
	var out strings.Builder

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
	var out strings.Builder
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

// ReturnStatement is a return value from a function.
type ReturnStatement struct {
	Token       lexer.Token // The "return" token
	ReturnValue Expression  // The expression value to return
}

func (rs *ReturnStatement) stmt()                {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out strings.Builder

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
