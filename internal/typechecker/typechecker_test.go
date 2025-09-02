package typechecker

import (
	"sigil/internal/lexer"
	"sigil/internal/parser"
	"testing"
)

func evalType(t *testing.T, input string) Type {
	t.Helper()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}

	tc := New()
	tp := tc.CheckProgram(program)

	if tc.HasErrors() {
		t.Fatalf("type checking errors: %v", tc.Errors())
	}

	return tp
}

func TestTypeCheckerBasics(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1 + 2", "Number"},
		{"let x: Number = 5;", "Void"},
		{"if (true) { 10 } else { 20 }", "Number"},
		{"if (false) { 10 } else { 20 }", "Number"},
		{"fun(x: Number): Number { x + 1 }", "(Number) -> Number"},
		{"fun(x: Number, y: Number): Number { x + y }", "(Number, Number) -> Number"},
		{"fun(): Number { 42 }", "() -> Number"},
	}

	for _, tt := range tests {
		got := evalType(t, tt.input)
		if got.String() != tt.want {
			t.Errorf("input %q: got %s, want %s", tt.input, got.String(), tt.want)
		}
	}
}

func TestTypeCheckerIfExpressions(t *testing.T) {
	tests := []struct {
		input       string
		shouldError bool
	}{
		{"if (true) { 10 } else { 20 }", false},
		{`if (false) { 10 } else { "string" }`, true}, // mismatched types
		{"if (42) { 10 } else { 20 }", true},          // non-bool condition
		{"if (true) { 10 }", false},                   // no else branch
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		tc := New()
		tc.CheckProgram(program)

		if tt.shouldError && !tc.HasErrors() {
			t.Errorf("expected errors for input %q, got none", tt.input)
		}
		if !tt.shouldError && tc.HasErrors() {
			t.Errorf("unexpected errors for input %q: %v", tt.input, tc.Errors())
		}
	}
}

func TestTypeCheckerFunctions(t *testing.T) {
	tests := []struct {
		input       string
		shouldError bool
	}{
		{"fun(x: Number): Number { x + 1 }(5)", false},
		{"fun(x: Number): Number { x + 1 }(\"hello\")", true}, // arg type mismatch
		{"fun(): Number { 42 }()", false},
		{"fun(x: Number): Number { let y: String = x; y }", true}, // assignment type mismatch
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		tc := New()
		tc.CheckProgram(program)

		if tt.shouldError && !tc.HasErrors() {
			t.Errorf("expected errors for input %q, got none", tt.input)
		}
		if !tt.shouldError && tc.HasErrors() {
			t.Errorf("unexpected errors for input %q: %v", tt.input, tc.Errors())
		}
	}
}

func TestTypeCheckerExtended(t *testing.T) {
	tests := []struct {
		input       string
		wantType    string
		shouldError bool
	}{
		// Arithmetic and prefix expressions
		{"-5", "Number", false},
		{"!true", "Boolean", false},
		{"-true", "", true}, // invalid unary minus
		{"!42", "", true},   // invalid logical not
		{"1 + 2 * 3", "Number", false},

		// Nested if expressions
		{"if (true) { if (false) { 1 } else { 2 } } else { 3 }", "Number", false},
		{"if (true) { 1 } else { if (false) { 2 } else { \"x\" } }", "", true}, // mismatched types

		// Variable shadowing
		{"let x: Number = 5; let x: Number = x + 1;", "Void", false},

		// Function returns
		{"fun(x: Number): Number { let y: Number = x + 1; y }", "(Number) -> Number", false},
		{"fun(x: Number): Number { let y: String = x; y }", "", true}, // type mismatch inside function

		// Function call with wrong argument
		{"fun(x: Number): Number { x + 1 }(\"wrong\")", "", true},

		// Equality comparisons
		{"1 == 2", "Boolean", false},
		{"1 == \"hi\"", "", true}, // invalid equality between different types

		// Boolean comparisons
		{"true == false", "Boolean", false},

		// Less than / greater than
		{"5 < 10", "Boolean", false},
		{"5 < \"x\"", "", true}, // invalid comparison
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		tc := New()
		got := tc.CheckProgram(program)

		if tt.shouldError && !tc.HasErrors() {
			t.Errorf("expected errors for input %q, got none", tt.input)
		}
		if !tt.shouldError && tc.HasErrors() {
			t.Errorf("unexpected errors for input %q: %v", tt.input, tc.Errors())
		}

		if !tt.shouldError && got.String() != tt.wantType {
			t.Errorf("input %q: got type %s, want %s", tt.input, got.String(), tt.wantType)
		}
	}
}
