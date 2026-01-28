package parser

import (
	"fmt"
	"regexp"
)

// A lexer.
type lexer struct {
	src      string
	pos      int
	tokens   []token
	patterns []pattern
}

// Represents a regex pattern for a token type.
type pattern struct {
	regex     *regexp.Regexp
	tokenType tokenType
}

// Returns whether the lexer is has reached EOF.
func (l *lexer) atEof() bool {
	return l.pos >= len(l.src)
}

// Get the remainder of the source.
func (l *lexer) remainder() string {
	return l.src[l.pos:]
}

// Advance the lexer pointer by n.
func (l *lexer) advance(n int) {
	l.pos += n
}

// Add a token of type t and value v to the tokens.
func (l *lexer) addToken(t tokenType, v string) {
	l.tokens = append(l.tokens, token{t, v})
}

// Creates a new lexer with all the patterns for lexing.
func newLexer(src string) *lexer {
	return &lexer{
		src:    src,
		tokens: make([]token, 0),
		patterns: []pattern{
			{regexp.MustCompile(`[0-9A-Za-z_]+`), TokenAlphaNum},
			{regexp.MustCompile(`=`), TokenEqual},
			{regexp.MustCompile(`\.`), TokenDot},
			{regexp.MustCompile(`\+`), TokenPlus},
			{regexp.MustCompile(`#`), TokenHash},
			{regexp.MustCompile(`-`), TokenDash},
			{regexp.MustCompile(`\*`), TokenStar},
			{regexp.MustCompile(`/`), TokenSlash},
			{regexp.MustCompile(`%`), TokenPercent},
			{regexp.MustCompile(`@`), TokenAt},
			{regexp.MustCompile(`\^`), TokenCaret},
			{regexp.MustCompile(`\(`), TokenLParen},
			{regexp.MustCompile(`\)`), TokenRParen},
			{regexp.MustCompile(`\[`), TokenLBracket},
			{regexp.MustCompile(`\]`), TokenRBracket},
			{regexp.MustCompile(`\{`), TokenLCurly},
			{regexp.MustCompile(`\}`), TokenRCurly},
			{regexp.MustCompile(`,`), TokenComma},
		},
	}
}

// Tokenize the source. Returns error for syntax error and when the source is empty.
func Tokenize(src string) ([]token, error) {
	l := newLexer(src)

	for !l.atEof() {
		var matched *pattern
		for l.remainder()[0] == ' ' {
			l.advance(1)
		}
		for _, p := range l.patterns {
			index := p.regex.FindStringIndex(l.remainder())
			if index != nil && index[0] == 0 {
				matched = &p
				break
			}
		}
		if matched == nil {
			return nil, fmt.Errorf("Invalid syntax: %s", l.remainder())
		}
		match := matched.regex.FindString(l.remainder())
		l.addToken(matched.tokenType, match)
		l.advance(len(match))
	}
	if len(l.tokens) == 0 {
		return nil, fmt.Errorf("Empty expression")
	}
	return l.tokens, nil
}
