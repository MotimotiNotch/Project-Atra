package lexer

import (
	"project-atra/token"
	"fmt"
)

type Lexer struct {
	input        string
	position     int  // Current position in input
	readPosition int  // Next position to read
	ch           byte // Current character under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
		l.position = len(l.input)
	} else {
		l.ch = l.input[l.readPosition]
		l.position = l.readPosition
	}
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) peekNChar(n int) byte {
	pos := l.position + n
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	if l.ch == 0 {
		tok.Literal = ""
		tok.Type = token.EOF
		return tok
	}

	if l.position == 0 || l.isAtLineStart() {
		indentCount := l.countIndent()
		if indentCount > 0 {
			tok.Type = token.INDENT
			tok.Literal = fmt.Sprintf("%d", indentCount)
			return tok
		}
	}

	l.skipHorizontalWhitespace()

	if l.ch == '#' {
		for l.ch != '\n' && l.ch != '\r' && l.ch != 0 {
			l.readChar()
		}
		return l.NextToken()
	}

	switch l.ch {
	case '\n':
		tok = token.Token{Type: token.NEWLINE, Literal: "\n"}
	case '\r':
		if l.peekChar() == '\n' {
			l.readChar()
			tok = token.Token{Type: token.NEWLINE, Literal: "\r\n"}
		} else {
			tok = token.Token{Type: token.NEWLINE, Literal: "\r"}
		}
	case '*':
		tok = newToken(token.STAR, l.ch)
	case '!':
		tok = newToken(token.SINGULARITY, l.ch)
	case '@':
		tok = newToken(token.REF, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '.':
		tok = newToken(token.DOT, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.FLOW, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '!' && l.peekNChar(2) == '-' {
			l.readChar() // read '!'
			l.readChar() // read '-'
			tok = token.Token{Type: token.REPEL, Literal: "-!-"}
		} else if isDigit(l.peekChar()) {
			tok.Literal = l.readNumber()
			tok.Type = token.NUMBER
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '+':
		if l.peekChar() == '+' && l.peekNChar(2) == '+' {
			l.readChar()
			l.readChar()
			tok = token.Token{Type: token.ATTRACT, Literal: "+++"}
		} else if isDigit(l.peekChar()) {
			tok.Literal = l.readNumber()
			tok.Type = token.NUMBER
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '~':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.WAVE_FLOW, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '=':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.WORMHOLE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.IDENT
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.NUMBER
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
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
	if l.ch == '+' || l.ch == '-' {
		l.readChar()
	}
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) isAtLineStart() bool {
	if l.position == 0 {
		return true
	}
	prevChar := l.input[l.position-1]
	return prevChar == '\n' || prevChar == '\r'
}

func (l *Lexer) countIndent() int {
	count := 0
	tempPos := l.position
	for tempPos < len(l.input) {
		ch := l.input[tempPos]
		if ch == ' ' {
			count++
		} else if ch == '\t' {
			count += 4
		} else {
			break
		}
		tempPos++
	}

	for i := 0; i < (tempPos - l.position); i++ {
		l.readChar()
	}

	return count
}

func (l *Lexer) skipHorizontalWhitespace() {
	for l.ch == ' ' || l.ch == '\t' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
