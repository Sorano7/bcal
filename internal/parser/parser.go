package parser

import (
	"fmt"
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

// If the next token is any of the arguments.
func (p *parser) nextNTokenIs(n int, ts ...tokenType) bool {
	if p.pos+n >= len(p.tokens) {
		return false
	}
	return slices.Contains(ts, p.tokens[p.pos+n].typeof)
}

// Parse a source string.
func Parse(src string) (Statement, error) {
	tokens, err := Tokenize(src)
	if err != nil {
		return nil, err
	}
	p := newParser(tokens)
	stmt := p.parseStmt()
	if isError(stmt) {
		return nil, fmt.Errorf("%s", stmt.(*ErrorNode).Msg)
	}
	return stmt, nil
}

// Parse a statement.
func (p *parser) parseStmt() Statement {
	switch {
	case p.currentTokenIs(TokenAt) && p.nextNTokenIs(1, TokenInt) && p.nextNTokenIs(2, TokenEqual):
		return p.parseAssignment()
	}
	stmt := &ExprStatement{}
	stmt.Expr = p.parseExpr(Lowest)
	if isError(stmt.Expr) {
		return stmt.Expr.(*ErrorNode)
	}
	return stmt
}

// Parse an assignment statement.
func (p *parser) parseAssignment() Statement {
	stmt := &AssignStatement{}
	stmt.Target = p.parseExpr(Lowest)
	if isError(stmt.Target) {
		return stmt.Target.(*ErrorNode)
	}
	if !p.currentTokenIs(TokenEqual) {
		return newErrorf("Missing equal sign in assignment to %s", stmt.Target)
	}
	p.advance()
	stmt.Value = p.parseExpr(Lowest)
	if isError(stmt.Value) {
		return stmt.Value.(*ErrorNode)
	}
	return stmt
}

// Parse an expression with the precedence.
func (p *parser) parseExpr(prec int) Expression {
	var expr Expression

	current := p.currentToken()
	switch current.typeof {
	case TokenInt, TokenDot, TokenLCurly:
		expr = p.parseNumber()
	case TokenDash:
		expr = p.parsePrefix()
	case TokenLParen:
		expr = p.parseGroup()
	case TokenLBracket:
		expr = p.parseBaseAnnotation()
	case TokenAt:
		expr = p.parseIdent()
	default:
		return newErrorf("Invalid syntax: %s", current.value)
	}
	if isError(expr) || p.atEof() {
		return expr
	}
	return p.handleInfix(prec, expr)
}

// Continue parsing infix for as long as the precedence allow.
func (p *parser) handleInfix(prec int, expr Expression) Expression {
	for !p.atEof() && prec < p.currentToken().prec() {
		switch {
		case p.currentTokenIs(TokenPlus, TokenDash, TokenStar, TokenSlash, TokenPercent, TokenCaret):
			return p.parseInfix(expr)
		case p.currentTokenIs(TokenArrow):
			return p.parseOutputBase(expr)
		}
	}
	return expr
}

// Parse a digit value.
func (p *parser) parseDigitValue() Expression {
	num := p.currentToken().value
	if !p.currentTokenIs(TokenInt) {
		return newErrorf("Not a int: %s", num)
	}
	n, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return newError(err)
	}
	p.advance()
	return &DigitValue{n}
}

// Parse a digit string.
func (p *parser) parseDigitString() Expression {
	if !p.currentTokenIs(TokenInt) {
		return newErrorf("Invalid digit string: %s", p.currentToken().value)
	}
	expr := &DigitString{Value: p.currentToken().value}
	p.advance()
	return expr
}

// Parse a list of digits.
func (p *parser) parseDigitList() Expression {
	p.advance()
	expr := &DigitList{}
	for !p.currentTokenIs(TokenRCurly) {
		next := p.parseDigitValue()
		if isError(next) {
			return next
		}
		expr.Value = append(expr.Value, next.(*DigitValue).Value)
		if p.currentTokenIs(TokenRCurly) {
			break
		}
		if p.currentTokenIs(TokenComma) {
			p.advance()
		}
	}
	p.advance()
	return expr
}

// Parse a digit string or list.
func (p *parser) parseStringOrList() Expression {
	if p.currentTokenIs(TokenLCurly) {
		return p.parseDigitList()
	}
	return p.parseDigitString()
}

// Parse a number.
func (p *parser) parseNumber() Expression {
	expr := &NumberLiteral{
		Int:    &DigitValue{0},
		Nonrep: &DigitValue{0},
		Rep:    &DigitValue{0},
	}
	if !p.currentTokenIs(TokenDot) {
		expr.Int = p.parseStringOrList()
		if isError(expr.Int) {
			return expr.Int
		}
	}
	return p.parseFracPart(expr)
}

// Parse the fraction part of a number.
func (p *parser) parseFracPart(expr *NumberLiteral) Expression {
	if !p.currentTokenIs(TokenDot) {
		return expr
	}
	p.advance()
	if !p.currentTokenIs(TokenLParen) {
		expr.Nonrep = p.parseStringOrList()
		if isError(expr.Nonrep) {
			return expr.Nonrep
		}
	}

	if p.currentTokenIs(TokenLParen) {
		p.advance()
		expr.Rep = p.parseStringOrList()
		if isError(expr.Rep) {
			return expr.Rep
		}
		if !p.currentTokenIs(TokenRParen) {
			return newError("Missing closing parenthesis")
		}
		p.advance()
	}
	return expr
}

// Parse base literal.
func (p *parser) parseBaseLiteral() Expression {
	p.advance()
	if !p.currentTokenIs(TokenInt) {
		return newError("Invalid base annotation: not an int")
	}
	lit := p.currentToken().value
	base, err := strconv.ParseInt(lit, 10, 64)
	if err != nil {
		return newErrorf("Invalid base annotation: %s", err)
	}
	p.advance()
	if !p.currentTokenIs(TokenRBracket) {
		return newError("Missing closing bracket")
	}
	p.advance()

	return &BaseAnnotation{Base: base}

}

// Parse a number with base annotation.
func (p *parser) parseBaseAnnotation() Expression {
	expr := p.parseBaseLiteral()
	if isError(expr) {
		return expr
	}
	e := expr.(*BaseAnnotation)
	e.Expr = p.parseExpr(Base)
	if isError(e.Expr) {
		return e.Expr
	}
	return e
}

// Parse prefix expressions.
func (p *parser) parsePrefix() Expression {
	expr := &PrefixExpr{Operator: p.currentToken().value}
	p.advance()
	expr.Right = p.parseExpr(Prefix)
	if isError(expr.Right) {
		return expr.Right
	}
	return expr
}

// Parse infix expressions.
func (p *parser) parseInfix(left Expression) Expression {
	expr := &InfixExpr{Left: left, Operator: p.currentToken().value}
	prec := p.currentToken().prec()
	p.advance()
	expr.Right = p.parseExpr(prec)
	if isError(expr.Right) {
		return expr.Right
	}
	return expr
}

// Parse group expressions.
func (p *parser) parseGroup() Expression {
	p.advance()
	expr := p.parseExpr(Lowest)
	if isError(expr) {
		return expr
	}
	if !p.currentTokenIs(TokenRParen) {
		return newError("Missing right parenthesis")
	}
	p.advance()
	return expr
}

// Parse an identifier.
func (p *parser) parseIdent() Expression {
	p.advance()
	expr := &Identifier{Name: ""}
	if !p.currentTokenIs(TokenInt) {
		return expr
	}
	expr.Name = p.currentToken().value
	p.advance()
	return expr
}

// Parse output expressions.
func (p *parser) parseOutputBase(left Expression) Expression {
	p.advance()
	base := p.parseBaseLiteral()
	if isError(base) {
		return base
	}
	expr := &OutputBase{Base: base.(*BaseAnnotation).Base, Expr: left}
	return expr
}
