package lexer

import (
	"unicode"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char
	line         int
	column       int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()

	// Capture position at the START of the token
	tokenLine := l.line
	tokenColumn := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    EQEQ,
				Literal: string(ch) + string(l.ch),
				Line:    tokenLine,
				Column:  tokenColumn,
			}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    ARROW,
				Literal: string(ch) + string(l.ch),
				Line:    tokenLine,
				Column:  tokenColumn,
			}
		} else {
			tok = l.newTokenAt(ASSIGN, string(l.ch), tokenLine, tokenColumn)
		}
	case '+':
		tok = l.newTokenAt(PLUS, string(l.ch), tokenLine, tokenColumn)
	case '-':
		tok = l.newTokenAt(MINUS, string(l.ch), tokenLine, tokenColumn)
	case '*':
		tok = l.newTokenAt(STAR, string(l.ch), tokenLine, tokenColumn)
	case '/':
		tok = l.newTokenAt(SLASH, string(l.ch), tokenLine, tokenColumn)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    NOTEQ,
				Literal: string(ch) + string(l.ch),
				Line:    tokenLine,
				Column:  tokenColumn,
			}
		} else {
			tok = l.newTokenAt(BANG, string(l.ch), tokenLine, tokenColumn)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    LTE,
				Literal: string(ch) + string(l.ch),
				Line:    tokenLine,
				Column:  tokenColumn,
			}
		} else {
			tok = l.newTokenAt(LT, string(l.ch), tokenLine, tokenColumn)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    GTE,
				Literal: string(ch) + string(l.ch),
				Line:    tokenLine,
				Column:  tokenColumn,
			}
		} else {
			tok = l.newTokenAt(GT, string(l.ch), tokenLine, tokenColumn)
		}
	case ';':
		tok = l.newTokenAt(SEMICOLON, string(l.ch), tokenLine, tokenColumn)
	case ',':
		tok = l.newTokenAt(COMMA, string(l.ch), tokenLine, tokenColumn)
	case ':':
		tok = l.newTokenAt(COLON, string(l.ch), tokenLine, tokenColumn)
	case '(':
		tok = l.newTokenAt(LPAREN, string(l.ch), tokenLine, tokenColumn)
	case ')':
		tok = l.newTokenAt(RPAREN, string(l.ch), tokenLine, tokenColumn)
	case '{':
		tok = l.newTokenAt(LBRACE, string(l.ch), tokenLine, tokenColumn)
	case '}':
		tok = l.newTokenAt(RBRACE, string(l.ch), tokenLine, tokenColumn)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		tok.Line = tokenLine
		tok.Column = tokenColumn
	case 0:
		tok.Type = EOF
		tok.Literal = ""
		tok.Line = tokenLine
		tok.Column = tokenColumn
	default:
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			tok.Type = LookupIdent(lit)
			tok.Literal = lit
			tok.Line = tokenLine
			tok.Column = tokenColumn
			return tok // Don't advance again - readIdentifier already did
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			tok.Line = tokenLine
			tok.Column = tokenColumn
			return tok // Don't advance again - readNumber already did
		} else {
			tok = l.newTokenAt(ILLEGAL, string(l.ch), tokenLine, tokenColumn)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) newToken(t TokenType, ch string) Token {
	return Token{
		Type:    t,
		Literal: ch,
		Line:    l.line,
		Column:  l.column,
	}
}

func (l *Lexer) newTokenAt(t TokenType, ch string, line, column int) Token {
	return Token{
		Type:    t,
		Literal: ch,
		Line:    line,
		Column:  column,
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	// TODO: handle decimal numbers later
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
