package parser

import (
	"sigil/internal/ast"
	"sigil/internal/lexer"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser had %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser Error: %s", msg)
	}
	t.FailNow()
}

func TestLetStatements(t *testing.T) {
	input := `
let x: Number = 5;
let y: Number = 10;
let foobar: Number = 123123;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Program.Statements does not contain 3 statements, got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.String not 'let', got=%q", stmt.String())
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt not *ast.LetStatement, got=%T", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s', got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.String() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s', got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 1234;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral() not 'return'. got=%q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have enough statements, got=%d, want=%d", len(program.Statements), 1)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected *ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected *ast.Identifier, got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value != %s, got %s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenValue() != %s, got %s", "foobar", ident.Value)
	}
}

func TestNumberLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have enough statements. Expected %d, got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("literal not *ast.NumberLiteral, got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value != 5, got=%f", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral != %s, got=%s", "5", literal.TokenLiteral())
	}
}
