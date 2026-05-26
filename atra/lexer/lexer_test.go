package lexer

import (
	"project-atra/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `*** Viscosity
Viscosity +++ Silent_Melody
`
	// Note: input newlines are returned as NEWLINE tokens
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.STAR, "*"},
		{token.STAR, "*"},
		{token.STAR, "*"},
		{token.IDENT, "Viscosity"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "Viscosity"},
		{token.ATTRACT, "+++"},
		{token.IDENT, "Silent_Melody"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIndent(t *testing.T) {
	input := `Viscosity:
  Silent
    *star
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "Viscosity"},
		{token.COLON, ":"},
		{token.NEWLINE, "\n"},
		{token.INDENT, "2"},
		{token.IDENT, "Silent"},
		{token.NEWLINE, "\n"},
		{token.INDENT, "4"},
		{token.STAR, "*"},
		{token.IDENT, "star"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
