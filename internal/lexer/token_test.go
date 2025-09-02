package lexer

import (
	"testing"
)

func TestLookupIdent(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"myVar", IDENT},
		{"let", LET},
		{"fun", FUNCTION},
		{"return", RETURN},
		{"true", TRUE},
		{"false", FALSE},
		{"if", IF},
		{"else", ELSE},
	}

	for _, tt := range tests {
		kw := LookupIdent(tt.input)
		if kw != tt.expected {
			t.Errorf("Wrong keyword returned. Expected=%s, Got=%s", tt.expected, kw)
		}
	}
}
