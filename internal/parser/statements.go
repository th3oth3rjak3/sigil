package parser

import (
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

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	name := p.curToken.Literal
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: name}

	// Optional type hint
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // consume ':'
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		stmt.TypeHint = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken() // move to expression
	stmt.Value = p.parseExpression(LOWEST)

	if fn, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fn.Name = name
		stmt.Value = fn
	}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

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
