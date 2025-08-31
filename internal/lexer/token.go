package lexer

type TokenType string

const (
	// Special Tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and Literals
	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	STRING = "STRING"

	// Keywords (some present now, more can be added later)
	LET      = "LET"
	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"

	// Single-char Operators
	ASSIGN = "ASSIGN"
	PLUS   = "PLUS"
	MINUS  = "MINUS"
	STAR   = "STAR"
	SLASH  = "SLASH"
	BANG   = "BANG"

	// Multi-char / comparison
	EQUAL                 = "EQUAL"
	NOT_EQUAL             = "NOT_EQUAL"
	LESS_THAN             = "LESS_THAN"
	LESS_THAN_OR_EQUAL    = "LESS_THAN_OR_EQUAL"
	GREATER_THAN          = "GREATER_THAN"
	GREATER_THAN_OR_EQUAL = "GREATER_THAN_OR_EQUAL"
	ARROW                 = "ARROW"

	// Delimiters
	COLON         = "COLON"
	SEMICOLON     = "SEMICOLON"
	COMMA         = "COMMA"
	LEFT_PAREN    = "LEFT_PAREN"
	RIGHT_PAREN   = "RIGHT_PAREN"
	LEFT_BRACE    = "LEFT_BRACE"
	RIGHT_BRACE   = "RIGHT_BRACE"
	LEFT_BRACKET  = "LEFT_BRACKET"
	RIGHT_BRACKET = "RIGHT_BRACKET"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]TokenType{
	"let":    LET,
	"fun":    FUNCTION,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
