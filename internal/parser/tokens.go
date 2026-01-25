package parser

type tokenType int

const (
	EOF tokenType = iota
	TokenInt
	TokenNonrep
	TokenRep
	TokenPlus
	TokenDash
	TokenStar
	TokenSlash
	TokenPercent
	TokenCaret
	TokenLParen
	TokenRParen
	TokenLBracket
	TokenRBracket
)

const (
	Lowest int = iota
	PlusMinus
	MulDiv
	Power
	Prefix
)

var opPrec = map[tokenType]int{
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
