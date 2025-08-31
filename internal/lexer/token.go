package lexer

type TokenType int

func (tt TokenType) String() string {
	switch tt {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"

	// Keywords
	case LET:
		return "LET"
	case FN:
		return "FN"
	case RETURN:
		return "RETURN"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"

	// Operators
	case ASSIGN:
		return "ASSIGN" // =
	case PLUS:
		return "PLUS" // +
	case MINUS:
		return "MINUS" // -
	case STAR:
		return "STAR" // *
	case SLASH:
		return "SLASH" // /
	case BANG:
		return "BANG" // !

	case EQEQ:
		return "EQEQ" // ==
	case NOTEQ:
		return "NOTEQ" // !=
	case LT:
		return "LT" // <
	case LTE:
		return "LTE" // <=
	case GT:
		return "GT" // >
	case GTE:
		return "GTE" // >=
	case ARROW:
		return "ARROW" // =>

	// Delimiters
	case COLON:
		return "COLON" // :
	case SEMICOLON:
		return "SEMICOLON" // ;
	case COMMA:
		return "COMMA" // ,
	case LPAREN:
		return "LPAREN" // (
	case RPAREN:
		return "RPAREN" // )
	case LBRACE:
		return "LBRACE" // {
	case RBRACE:
		return "RBRACE" // }

	default:
		return "UNKNOWN"
	}
}

const (
	// Special Tokens
	ILLEGAL TokenType = iota
	EOF

	// Identifiers and Literals
	IDENT
	NUMBER
	STRING

	// Keywords (some present now, more can be added later)
	LET
	PRINT
	FN
	RETURN
	TRUE
	FALSE

	// Single-char Operators
	ASSIGN // =
	PLUS   // +
	MINUS  // -
	STAR   // *
	SLASH  // /
	BANG   // !

	// Multi-char / comparison
	EQEQ  // ==
	NOTEQ // !=
	LT    // <
	LTE   // <=
	GT    // >
	GTE   // >=
	ARROW // =>

	// Delimiters
	COLON
	SEMICOLON
	COMMA
	LPAREN
	RPAREN
	LBRACE
	RBRACE
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]TokenType{
	"let":    LET,
	"print":  PRINT,
	"fn":     FN,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
