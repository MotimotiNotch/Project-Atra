package parser

import (
	"project-atra/ast"
	"project-atra/lexer"
	"project-atra/token"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	FLOW  // ->, ~>, =>
	FORCE // +++, -!-
	INDEX // :
	CALL  // (
)

var precedences = map[token.TokenType]int{
	token.FLOW:      FLOW,
	token.WAVE_FLOW: FLOW,
	token.WORMHOLE:  FLOW,
	token.ATTRACT:   FORCE,
	token.REPEL:     FORCE,
	token.COLON:     INDEX,
	token.LBRACE:    INDEX,
	token.LPAREN:    CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	currentIndent int
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseCelestialNode)
	p.registerPrefix(token.STAR, p.parseCelestialNode)
	p.registerPrefix(token.SINGULARITY, p.parseCelestialNode)
	p.registerPrefix(token.REF, p.parseRhizomeExpression)
	p.registerPrefix(token.NUMBER, p.parseNumberLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.FLOW, p.parseFlowStatement)
	p.registerInfix(token.WAVE_FLOW, p.parseFlowStatement)
	p.registerInfix(token.WORMHOLE, p.parseFlowStatement)
	p.registerInfix(token.ATTRACT, p.parseForceStatement)
	p.registerInfix(token.REPEL, p.parseForceStatement)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.COLON, p.parseKeyValuePair)
	p.registerInfix(token.LBRACE, p.parseBlockExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseUniverse() *ast.Universe {
	universe := &ast.Universe{}
	universe.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		if p.curToken.Type == token.NEWLINE || p.curToken.Type == token.INDENT {
			p.nextToken()
			continue
		}
		stmt := p.parseStatement()
		if stmt != nil {
			universe.Statements = append(universe.Statements, stmt)
		}
		p.nextToken()
	}

	return universe
}

func (p *Parser) parseStatement() ast.Statement {
	return p.parseExpressionStatement()
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseCelestialNode() ast.Expression {
	node := &ast.CelestialNode{Token: p.curToken}

	for p.curToken.Type == token.STAR || p.curToken.Type == token.SINGULARITY {
		if p.curToken.Type == token.STAR {
			node.Mass++
		} else {
			node.IsSingularity = true
		}
		p.nextToken()
	}

	if p.curToken.Type != token.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("expected identifier, got %s", p.curToken.Type))
		return nil
	}

	node.Value = p.curToken.Literal

	// Infix calls like Phase(Name) will handle the LBRACE after parsing CallExpression.
	// But simple cases like *star: or Phase: need handling here.
	if p.peekTokenIs(token.COLON) || p.peekTokenIs(token.LBRACE) {
		// Wait, if it's an IDENT followed by LBRACE, parseExpression will handle LBRACE as infix.
		// If it's a COLON, it's also infix.
	}

	return node
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	if p.curTokenIs(token.LBRACE) {
		p.nextToken() // Skip {
		for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
			if p.curTokenIs(token.NEWLINE) || p.curTokenIs(token.INDENT) {
				p.nextToken()
				continue
			}
			stmt := p.parseStatement()
			if stmt != nil {
				block.Statements = append(block.Statements, stmt)
			}
			p.nextToken()
		}
		return block
	}

	// Legacy Indentation Logic (Colon-based)
	for !p.curTokenIs(token.NEWLINE) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	p.nextToken() // Skip NEWLINE
	if !p.curTokenIs(token.INDENT) {
		return block
	}

	indent, _ := strconv.Atoi(p.curToken.Literal)
	if indent <= p.currentIndent {
		return block
	}

	oldIndent := p.currentIndent
	p.currentIndent = indent
	p.nextToken() // Skip INDENT

	for !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.NEWLINE) {
			p.nextToken()
			if p.curTokenIs(token.INDENT) {
				newIndent, _ := strconv.Atoi(p.curToken.Literal)
				if newIndent < p.currentIndent {
					p.currentIndent = oldIndent
					return block
				}
				p.nextToken()
			} else {
				p.currentIndent = oldIndent
				return block
			}
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	p.currentIndent = oldIndent
	return block
}

func (p *Parser) parseRhizomeExpression() ast.Expression {
	exp := &ast.RhizomeExpression{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	exp.Target = p.curToken.Literal
	return exp
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	lit := &ast.NumberLiteral{Token: p.curToken}
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as float64", p.curToken.Literal))
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseKeyValuePair(left ast.Expression) ast.Expression {
	kv := &ast.KeyValuePair{Token: p.curToken, Key: left}
	p.nextToken()
	kv.Value = p.parseExpression(LOWEST)
	return kv
}

func (p *Parser) parseBlockExpression(left ast.Expression) ast.Expression {
	if node, ok := left.(*ast.CelestialNode); ok {
		node.Body = p.parseBlockStatement()
		return node
	}
	if call, ok := left.(*ast.CallExpression); ok {
		// If it's a CallExpression, we need to wrap it or attach the body.
		// For simplicity, let's treat it as a block attached to the call result.
		p.parseBlockStatement()
		return call
	}
	p.parseBlockStatement()
	return left
}

func (p *Parser) parseFlowStatement(left ast.Expression) ast.Expression {
	stmt := &ast.FlowStatement{Token: p.curToken, Operator: p.curToken.Literal, Source: left}
	precedence := p.curPrecedence()
	p.nextToken()
	stmt.Target = p.parseExpression(precedence)
	return stmt
}

func (p *Parser) parseForceStatement(left ast.Expression) ast.Expression {
	stmt := &ast.ForceStatement{Token: p.curToken, Operator: p.curToken.Literal, Left: left}
	precedence := p.curPrecedence()
	p.nextToken()
	stmt.Right = p.parseExpression(precedence)
	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool { return p.curToken.Type == t }
func (p *Parser) peekTokenIs(t token.TokenType) bool { return p.peekToken.Type == t }
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}
func (p *Parser) Errors() []string { return p.errors }
func (p *Parser) peekError(t token.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
}
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s found", t))
}
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) { p.prefixParseFns[tokenType] = fn }
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) { p.infixParseFns[tokenType] = fn }
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok { return p }
	return LOWEST
}
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok { return p }
	return LOWEST
}
