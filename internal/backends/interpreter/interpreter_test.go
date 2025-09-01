package interpreter

import (
	"sigil/internal/lexer"
	"sigil/internal/parser"
	"testing"
)

// evalInput executes a program and returns the value of the last non-void expression.
func evalInput(t *testing.T, input string) Value {
	t.Helper()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}

	interp := New().(*Interpreter)

	var result Value
	for _, stmt := range program.Statements {
		val, err := interp.ExecuteStatement(stmt)
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
		if val != nil {
			result = val
		}
	}

	return result
}

func TestInterpreterExecution(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{"1 + 2", 3.0},
		{"let x: Number = 5; x", 5.0},
		{"if (true) { 10 } else { 20 }", 10.0},
		{"if (false) { 10 } else { 20 }", 20.0},
		{"fun(x: Number): Number { x + 1 }(5)", 6.0},
		{"let add = fun(x: Number, y: Number): Number { x + y }; add(2,3)", 5.0},
		{`"Hello" + " World!"`, "Hello World!"},
		// String concatenation with conversion
		// {`"Hello " + string(42)`, "Hello 42"}, // Can't do this yet, no builtin string function
		// Edge case: concatenating non-strings should produce an error
		{`"Hello " + 42`, "error"},
		{`1 + " World!"`, "error"},
	}

	for _, tt := range tests {
		if tt.want == "error" {
			// Expecting an execution error
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			if len(p.Errors()) > 0 {
				t.Fatalf("parse errors: %v", p.Errors())
			}

			interp := New().(*Interpreter)
			var execErr error
			for _, stmt := range program.Statements {
				_, execErr = interp.ExecuteStatement(stmt)
				if execErr != nil {
					break
				}
			}

			if execErr == nil {
				t.Errorf("input %q: expected execution error, got none", tt.input)
			}

			continue
		}

		got := evalInput(t, tt.input)

		switch want := tt.want.(type) {
		case float64:
			num, ok := got.(*NumberValue)
			if !ok {
				t.Fatalf("input %q: expected NumberValue, got %T", tt.input, got)
			}
			if num.Value != want {
				t.Errorf("input %q: got %v, want %v", tt.input, num.Value, want)
			}
		case bool:
			b, ok := got.(*BoolValue)
			if !ok {
				t.Fatalf("input %q: expected BoolValue, got %T", tt.input, got)
			}
			if b.Value != want {
				t.Errorf("input %q: got %v, want %v", tt.input, b.Value, want)
			}
		case string:
			s, ok := got.(*StringValue)
			if !ok {
				t.Fatalf("input %q: expected StringValue, got %T", tt.input, got)
			}
			if s.Value != want {
				t.Errorf("input %q: got %q, want %q", tt.input, s.Value, want)
			}
		default:
			t.Fatalf("unsupported test type %T", tt.want)
		}
	}
}
