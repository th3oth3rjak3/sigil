package parser

import (
	"fmt"
	"sigil/internal/ast"
	"sigil/internal/lexer"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	fmt.Printf("[parseLetStatement] start: curToken=%s, peekToken=%s\n", p.curToken.Literal, p.peekToken.Literal)
	stmt := &ast.LetStatement{Token: p.curToken}

	// Expect identifier after 'let'
	if !p.expectPeek(lexer.IDENT) {
		fmt.Println("[parseLetStatement] expected IDENT after 'let'")
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	fmt.Printf("[parseLetStatement] parsed name: %s\n", stmt.Name.Value)

	// Optional type annotation
	if p.peekTokenIs(lexer.COLON) {
		fmt.Println("[parseLetStatement] detected type hint")
		p.nextToken() // consume ':'

		if !p.expectPeek(lexer.IDENT) {
			fmt.Println("[parseLetStatement] expected type name after ':'")
			return nil
		}

		stmt.TypeHint = p.parseType()
		if stmt.TypeHint == nil {
			fmt.Println("[parseLetStatement] failed to parse type hint")
			return nil
		}
		fmt.Printf("[parseLetStatement] parsed type hint: %+v\n", stmt.TypeHint)

	}

	// After optional type, expect '='
	if !p.expectPeek(lexer.ASSIGN) {
		fmt.Printf("[parseLetStatement] expected '=', got %s\n", p.peekToken.Literal)
		return nil
	}
	p.nextToken() // move to expression
	fmt.Printf("[parseLetStatement] parsing value expression at token: %s\n", p.curToken.Literal)

	stmt.Value = p.parseExpression(LOWEST)
	if stmt.Value == nil {
		fmt.Println("[parseLetStatement] failed to parse value expression")
		return nil
	}

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
		fmt.Println("[parseLetStatement] consumed semicolon")
	}

	fmt.Printf("[parseLetStatement] finished: curToken=%s, peekToken=%s\n", p.curToken.Literal, p.peekToken.Literal)
	return stmt
}

// func (p *Parser) parseLetStatement() ast.Statement {
// 	fmt.Printf("[parseLetStatement] curToken=%s, peekToken=%s\n", p.curToken.Literal, p.peekToken.Literal)
// 	stmt := &ast.LetStatement{Token: p.curToken}

// 	// Expect an identifier
// 	if !p.expectPeek(lexer.IDENT) {
// 		fmt.Println("[parseLetStatement] expected IDENT after 'let'")
// 		return nil
// 	}

// 	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
// 	fmt.Printf("[parseLetStatement] parsed name: %s\n", stmt.Name.Value)

// 	// Optional type annotation
// 	if p.peekTokenIs(lexer.COLON) {
// 		fmt.Println("[parseLetStatement] detected type hint")
// 		p.nextToken() // consume ':'
// 		p.nextToken() // move to type token
// 		stmt.TypeHint = p.parseType()
// 		if stmt.TypeHint == nil {
// 			fmt.Println("[parseLetStatement] failed to parse type hint")
// 			return nil
// 		}
// 		fmt.Printf("[parseLetStatement] parsed type hint: %+v\n", stmt.TypeHint)
// 	}

// 	// After optional type, expect '='
// 	if !p.expectPeek(lexer.ASSIGN) {
// 		fmt.Printf("[parseLetStatement] expected '=', got %s\n", p.peekToken.Literal)
// 		return nil
// 	}

// 	p.nextToken() // move to expression
// 	fmt.Printf("[parseLetStatement] parsing value expression at token: %s\n", p.curToken.Literal)
// 	stmt.Value = p.parseExpression(LOWEST)
// 	if stmt.Value == nil {
// 		fmt.Println("[parseLetStatement] failed to parse value expression")
// 		return nil
// 	}

// 	if p.peekTokenIs(lexer.SEMICOLON) {
// 		p.nextToken() // consume ';'
// 		fmt.Println("[parseLetStatement] consumed semicolon")
// 	}

// 	return stmt
// }

// func (p *Parser) parseLetStatement() ast.Statement {
// 	stmt := &ast.LetStatement{Token: p.curToken}

// 	if !p.expectPeek(lexer.IDENT) {
// 		return nil
// 	}

// 	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

// 	// Optional type annotation
// 	if p.peekTokenIs(lexer.COLON) {
// 		p.nextToken() // consume ':'
// 		p.nextToken() // move to type token
// 		stmt.TypeHint = p.parseType()
// 		if stmt.TypeHint == nil {
// 			return nil
// 		}
// 	}

// 	// After optional type, advance to '='
// 	if !p.expectPeek(lexer.ASSIGN) {
// 		return nil
// 	}
// 	p.nextToken() // advance to start of value expression

// 	stmt.Value = p.parseExpression(LOWEST)
// 	if stmt.Value == nil {
// 		return nil
// 	}

// 	if p.peekTokenIs(lexer.SEMICOLON) {
// 		p.nextToken() // consume ';'
// 	}

// 	return stmt
// }

// func (p *Parser) parseLetStatement() *ast.LetStatement {
// 	stmt := &ast.LetStatement{Token: p.curToken}

// 	if !p.expectPeek(lexer.IDENT) {
// 		return nil
// 	}

// 	name := p.curToken.Literal
// 	stmt.Name = &ast.Identifier{Token: p.curToken, Value: name}

// 	// Optional type hint
// 	if p.peekTokenIs(lexer.COLON) {
// 		p.nextToken() // consume ':'
// 		if !p.expectPeek(lexer.IDENT) {
// 			return nil
// 		}
// 		stmt.TypeHint = p.parseType()
// 	}

// 	if !p.expectPeek(lexer.ASSIGN) {
// 		return nil
// 	}

// 	p.nextToken() // move to expression
// 	stmt.Value = p.parseExpression(LOWEST)

// 	if fn, ok := stmt.Value.(*ast.FunctionLiteral); ok {
// 		fn.Name = name
// 		stmt.Value = fn
// 	}

// 	if p.peekTokenIs(lexer.SEMICOLON) {
// 		p.nextToken()
// 	}

// 	return stmt
// }

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken() // move past 'return'

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
		stmt.HasSemicolon = true
	} else {
		stmt.HasSemicolon = false
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.nextToken()

	for !p.curTokenIs(lexer.RIGHT_BRACE) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}

	return block
}
