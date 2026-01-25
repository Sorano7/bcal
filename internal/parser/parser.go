package parser

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
)

// A parser for constructing AST.
type parser struct {
	tokens []token
	pos    int
}

// Create a new parser with a slice of tokens.
func newParser(tokens []token) *parser {
	return &parser{
		tokens: tokens,
	}
}

// Advance the parser position by one.
func (p *parser) advance() {
	p.pos++
}

// Return whether the parser is at eof.
func (p *parser) atEof() bool {
	return p.pos >= len(p.tokens)
}

// Return the current token.
func (p *parser) currentToken() token {
	if p.atEof() {
		return token{EOF, ""}
	}
	return p.tokens[p.pos]
}

// If the current token is any of the arguments.
func (p *parser) currentTokenIs(ts ...tokenType) bool {
	return slices.Contains(ts, p.currentToken().typeof)
}

// Parse a source string.
func Parse(src string) (Expression, error) {
	tokens, err := Tokenize(src)
	if err != nil {
		return nil, err
	}
	p := newParser(tokens)
	return p.parseExpr(Lowest)
}

// Parse an expression with the precedence.
func (p *parser) parseExpr(prec int) (Expression, error) {
	var expr Expression
	var err error

	current := p.currentToken()
	switch current.typeof {
	case TokenInt:
		expr, err = p.parseNumber()
	case TokenDash:
		expr, err = p.parsePrefix()
	case TokenLParen:
		expr, err = p.parseGroup()
	case TokenLBracket:
		expr, err = p.parseBaseNumber()
	}
	if err != nil {
		return nil, err
	}
	if p.atEof() {
		return expr, nil
	}
	return p.handleInfix(prec, expr)
}

// Continue parsing infix for as long as the precedence allow.
func (p *parser) handleInfix(prec int, expr Expression) (Expression, error) {
	var err error
	for !p.atEof() && prec < p.currentToken().prec() {
		switch {
		case p.currentTokenIs(TokenPlus, TokenDash, TokenStar, TokenSlash, TokenPercent, TokenCaret):
			expr, err = p.parseInfix(expr)
		default:
			return expr, nil
		}
	}
	if err != nil {
		return nil, err
	}
	return expr, nil
}

// Parse the fraction part of a number.
func (p *parser) parseFracPart(expr *NumberLiteral) (Expression, error) {
	if p.currentTokenIs(TokenRep) {
		clean := regexp.MustCompile(`[0-9A-Za-z]+`).FindString(p.currentToken().value)
		expr.Rep = clean
		p.advance()
		return expr, nil
	}

	if p.currentTokenIs(TokenNonrep) {
		noDot := p.currentToken().value[1:]
		expr.Nonrep = noDot
		p.advance()
	}

	if p.currentTokenIs(TokenLParen) {
		p.advance()
		if !p.currentTokenIs(TokenInt) {
			return nil, fmt.Errorf("Period must be a series of digits")
		}
		expr.Rep = p.currentToken().value
		p.advance()
		if !p.currentTokenIs(TokenRParen) {
			return nil, fmt.Errorf("Missing closing parenthesis")
		}
		p.advance()
	}
	return expr, nil
}

// Parse a number.
func (p *parser) parseNumber() (Expression, error) {
	expr := &NumberLiteral{Int: "0"}
	if p.currentTokenIs(TokenInt) {
		expr.Int = p.currentToken().value
		p.advance()
	}
	return p.parseFracPart(expr)
}

// Parse a number with base annotation.
func (p *parser) parseBaseNumber() (Expression, error) {
	p.advance()
	if !p.currentTokenIs(TokenInt) {
		return nil, fmt.Errorf("Invalid base annotation: not an int")
	}
	lit := p.currentToken().value
	base, err := strconv.ParseInt(lit, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid base annotation: %s", err)
	}
	p.advance()
	if !p.currentTokenIs(TokenRBracket) {
		return nil, fmt.Errorf("Missing closing bracket")
	}
	p.advance()
	expr, err := p.parseNumber()
	if err != nil {
		return nil, err
	}
	expr.(*NumberLiteral).Base = int(base)
	return expr, nil
}

// Parse prefix expressions.
func (p *parser) parsePrefix() (Expression, error) {
	expr := &PrefixExpr{Operator: p.currentToken().value}
	p.advance()
	right, err := p.parseExpr(Prefix)
	if err != nil {
		return nil, err
	}
	expr.Right = right
	return expr, nil
}

// Parse infix expressions.
func (p *parser) parseInfix(left Expression) (Expression, error) {
	expr := &InfixExpr{Left: left, Operator: p.currentToken().value}
	prec := p.currentToken().prec()
	p.advance()
	right, err := p.parseExpr(prec)
	if err != nil {
		return nil, err
	}
	expr.Right = right
	return expr, nil
}

// Parse group expressions.
func (p *parser) parseGroup() (Expression, error) {
	p.advance()
	expr, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}
	if !p.currentTokenIs(TokenRParen) {
		return nil, fmt.Errorf("Missing right parenthesis")
	}
	p.advance()
	return expr, nil
}
