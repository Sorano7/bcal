package parser

type Expression interface {
	expr()
}

type NumberLiteral struct {
	Int    string
	Nonrep string
	Rep    string
	Base   int
}

func (d *NumberLiteral) expr() {}

type InfixExpr struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpr) expr() {}

type PrefixExpr struct {
	Operator string
	Right    Expression
}

func (p *PrefixExpr) expr() {}
