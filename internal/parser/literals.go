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

func (p *Parser) parseType() ast.Type {
	switch p.curToken.Type {
	case lexer.IDENT:
		t := &ast.SimpleType{Token: p.curToken, Name: p.curToken.Literal}
		// Remove this line: p.nextToken()
		return t

	case lexer.LEFT_PAREN:
		p.nextToken() // consume '('
		paramTypes := []ast.Type{}

		for !p.curTokenIs(lexer.RIGHT_PAREN) && !p.curTokenIs(lexer.EOF) {
			paramTypes = append(paramTypes, p.parseType())
			p.nextToken() // Advance past the parsed type
			if p.curTokenIs(lexer.COMMA) {
				p.nextToken()
			}
		}

		if !p.curTokenIs(lexer.RIGHT_PAREN) {
			p.errors = append(p.errors, "expected RIGHT_PAREN at end of function type parameter list")
			return nil
		}

		p.nextToken() // consume RIGHT_PAREN

		if !p.curTokenIs(lexer.ARROW) {
			p.errors = append(p.errors, "expected ARROW after function type parameter list")
			return nil
		}

		p.nextToken() // consume ARROW
		returnType := p.parseType()
		if returnType == nil {
			return nil
		}
		// Don't advance here - let the caller handle it

		return &ast.FunctionType{
			ParamTypes: paramTypes,
			ReturnType: returnType,
		}
	}

	p.errors = append(p.errors, fmt.Sprintf("unexpected token in type: %s", p.curToken.Literal))
	return nil
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(lexer.LEFT_PAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.curTokenIs(lexer.COLON) {
		p.errors = append(p.errors, "expected ':' before return type")
		return nil
	}
	p.nextToken() // move to the return type token

	lit.ReturnType = p.parseType()
	if lit.ReturnType == nil {
		return nil
	}

	// NEW: Advance past the return type since parseType() no longer advances
	p.nextToken()

	if !p.curTokenIs(lexer.LEFT_BRACE) {
		p.errors = append(p.errors, fmt.Sprintf("expected '{' after function literal, got %s", p.curToken.Literal))
		return nil
	}

	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.FunctionParameter {
	fmt.Printf("[parseFunctionParameters] curToken=%s, peekToken=%s\n", p.curToken.Literal, p.peekToken.Literal)
	parameters := []*ast.FunctionParameter{}

	// Empty parameter shortcut: fun()
	if p.peekTokenIs(lexer.RIGHT_PAREN) {
		p.nextToken() // move to ')'
		p.nextToken() // consume ')'
		fmt.Println("[parseFunctionParameters] empty parameter list")
		return parameters
	}

	// Move into first parameter
	p.nextToken() // advance to first parameter
	for {
		fmt.Printf("[parseFunctionParameters] parsing parameter: curToken=%s\n", p.curToken.Literal)
		name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if !p.expectPeek(lexer.COLON) {
			fmt.Printf("[parseFunctionParameters] expected ':', got %s\n", p.peekToken.Literal)
			return nil
		}

		p.nextToken() // move to type
		typeHint := p.parseType()
		if typeHint == nil {
			fmt.Printf("[parseFunctionParameters] failed to parse type for %s\n", name.Value)
			return nil
		}

		// NEW: Advance past the type token since parseType() no longer does this
		p.nextToken()

		parameters = append(parameters, &ast.FunctionParameter{
			Name:     name,
			TypeHint: typeHint,
		})
		fmt.Printf("[parseFunctionParameters] added parameter: %s : %+v\n", name.Value, typeHint)

		// Move to next parameter or finish
		if p.curTokenIs(lexer.COMMA) {
			p.nextToken() // move to next parameter name
			continue
		}

		if p.curTokenIs(lexer.RIGHT_PAREN) {
			break
		}

		// Syntax error if unexpected token
		p.errors = append(p.errors, fmt.Sprintf("expected ',' or ')', got %s", p.curToken.Literal))
		fmt.Printf("[parseFunctionParameters] unexpected token %s\n", p.curToken.Literal)
		return nil
	}

	p.nextToken() // consume RIGHT_PAREN
	fmt.Println("[parseFunctionParameters] finished parsing parameters")
	return parameters
}
