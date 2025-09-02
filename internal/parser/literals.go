package parser

import (
	"fmt"
	"sigil/internal/ast"
	"sigil/internal/lexer"
	"strconv"
)

func (p *Parser) parseNumberLiteral() ast.Expression {
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("Error parsing Number: %s", err)
		p.errors = append(p.errors, msg)
	}
	return &ast.NumberLiteral{Token: p.curToken, Value: value}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(lexer.TRUE)}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(lexer.LEFT_PAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	lit.ReturnType = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.LEFT_BRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.FunctionParameter {
	parameters := []*ast.FunctionParameter{}

	// Shortcut for empty parameter list
	if p.peekTokenIs(lexer.RIGHT_PAREN) {
		p.nextToken() // consume ')'
		return parameters
	}

	for {
		p.nextToken() // move to parameter name
		name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if !p.expectPeek(lexer.COLON) {
			return nil
		}

		p.nextToken() // move to type
		typeHint := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		parameters = append(parameters, &ast.FunctionParameter{Name: name, TypeHint: typeHint})

		if !p.peekTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken() // consume comma
	}

	if !p.expectPeek(lexer.RIGHT_PAREN) {
		return nil
	}

	return parameters
}
