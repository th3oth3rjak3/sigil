package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five: Number = 5;
let ten: Number = 10;
let add = fun(x: Number, y: Number): Number => x + y;
let result: Number = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LET, "let"},
		{IDENT, "five"},
		{COLON, ":"},
		{IDENT, "Number"},
		{ASSIGN, "="},
		{NUMBER, "5"},
		{SEMICOLON, ";"},

		{LET, "let"},
		{IDENT, "ten"},
		{COLON, ":"},
		{IDENT, "Number"},
		{ASSIGN, "="},
		{NUMBER, "10"},
		{SEMICOLON, ";"},

		{LET, "let"},
		{IDENT, "add"},
		{ASSIGN, "="},
		{FUNCTION, "fun"},
		{LEFT_PAREN, "("},
		{IDENT, "x"},
		{COLON, ":"},
		{IDENT, "Number"},
		{COMMA, ","},
		{IDENT, "y"},
		{COLON, ":"},
		{IDENT, "Number"},
		{RIGHT_PAREN, ")"},
		{COLON, ":"},
		{IDENT, "Number"},
		{ARROW, "=>"},
		{IDENT, "x"},
		{PLUS, "+"},
		{IDENT, "y"},
		{SEMICOLON, ";"},

		{LET, "let"},
		{IDENT, "result"},
		{COLON, ":"},
		{IDENT, "Number"},
		{ASSIGN, "="},
		{IDENT, "add"},
		{LEFT_PAREN, "("},
		{IDENT, "five"},
		{COMMA, ","},
		{IDENT, "ten"},
		{RIGHT_PAREN, ")"},
		{SEMICOLON, ";"},

		{BANG, "!"},
		{MINUS, "-"},
		{SLASH, "/"},
		{STAR, "*"},
		{NUMBER, "5"},
		{SEMICOLON, ";"},

		{NUMBER, "5"},
		{LESS_THAN, "<"},
		{NUMBER, "10"},
		{GREATER_THAN, ">"},
		{NUMBER, "5"},
		{SEMICOLON, ";"},

		{IF, "if"},
		{LEFT_PAREN, "("},
		{NUMBER, "5"},
		{LESS_THAN, "<"},
		{NUMBER, "10"},
		{RIGHT_PAREN, ")"},
		{LEFT_BRACE, "{"},
		{RETURN, "return"},
		{TRUE, "true"},
		{SEMICOLON, ";"},
		{RIGHT_BRACE, "}"},
		{ELSE, "else"},
		{LEFT_BRACE, "{"},
		{RETURN, "return"},
		{FALSE, "false"},
		{SEMICOLON, ";"},
		{RIGHT_BRACE, "}"},

		{NUMBER, "10"},
		{EQUAL, "=="},
		{NUMBER, "10"},
		{SEMICOLON, ";"},

		{NUMBER, "10"},
		{NOT_EQUAL, "!="},
		{NUMBER, "9"},
		{SEMICOLON, ";"},

		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - TokenType wrong, have %s, want %s", i, tok.Type, tt.expectedType)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal Value wrong, have %s, want %s", i, tok.Literal, tt.expectedLiteral)
		}
	}
}
