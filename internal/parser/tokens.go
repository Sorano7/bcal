package parser

type tokenType int

const (
	EOF tokenType = iota
	TokenAlphaNum
	TokenEqual
	TokenDot
	TokenPlus
	TokenDash
	TokenStar
	TokenSlash
	TokenPercent
	TokenCaret
	TokenAt
	TokenLParen
	TokenRParen
	TokenLBracket
	TokenRBracket
	TokenLCurly
	TokenRCurly
	TokenHash
	TokenComma
)

const (
	Lowest int = iota
	Output
	Base
	PlusMinus
	MulDiv
	Power
	Prefix
)

var opPrec = map[tokenType]int{
	TokenHash:    Output,
	TokenPlus:    PlusMinus,
	TokenDash:    PlusMinus,
	TokenStar:    MulDiv,
	TokenSlash:   MulDiv,
	TokenPercent: MulDiv,
	TokenCaret:   Power,
}

type token struct {
	typeof tokenType
	value  string
}

func (t token) prec() int {
	prec, ok := opPrec[t.typeof]
	if !ok {
		prec = Lowest
	}
	return prec
}
