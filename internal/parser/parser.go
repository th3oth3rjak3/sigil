package parser

import (
	"fmt"
	"sigil/internal/ast"
	"sigil/internal/lexer"
)

// Precedences
const (
	_ int = iota
	LOWEST
	EQUALS       // ==
	LESS_GREATER // < or >
	SUM          // + -
	PRODUCT      // * /
	PREFIX       // -x, !x
	CALL         // myFunction(x)
)

var precedences = map[lexer.TokenType]int{
	lexer.EQUAL:                 EQUALS,
	lexer.NOT_EQUAL:             EQUALS,
	lexer.LESS_THAN:             LESS_GREATER,
	lexer.LESS_THAN_OR_EQUAL:    LESS_GREATER,
	lexer.GREATER_THAN:          LESS_GREATER,
	lexer.GREATER_THAN_OR_EQUAL: LESS_GREATER,
	lexer.PLUS:                  SUM,
	lexer.MINUS:                 SUM,
	lexer.STAR:                  PRODUCT,
	lexer.SLASH:                 PRODUCT,
	lexer.LEFT_PAREN:            CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read the first two tokens to preload the data in the parser
	// otherwise the tokens will be empty.
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.LEFT_PAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.BANG, p.parsePrefixExpression)
	p.registerPrefix(lexer.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.IF, p.parseIfExpression)
	p.registerPrefix(lexer.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.STAR, p.parseInfixExpression)
	p.registerInfix(lexer.SLASH, p.parseInfixExpression)
	p.registerInfix(lexer.EQUAL, p.parseInfixExpression)
	p.registerInfix(lexer.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(lexer.GREATER_THAN, p.parseInfixExpression)
	p.registerInfix(lexer.LESS_THAN, p.parseInfixExpression)
	p.registerInfix(lexer.GREATER_THAN_OR_EQUAL, p.parseInfixExpression)
	p.registerInfix(lexer.LESS_THAN_OR_EQUAL, p.parseInfixExpression)
	p.registerInfix(lexer.LEFT_PAREN, p.parseCallExpression)

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}

	return program
}

// Precedence helpers
func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

// Registration
func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Helpers
func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
