package mathb

import (
	"fmt"
	"math/big"
)

// Represents a base rational in normalized form.
type Rational struct {
	num   *big.Int
	denom *big.Int
	base  int64
}

// Sets n's base to base.
func (n *Rational) WithBase(base int64) *Rational {
	n.base = base
	return n
}

// Allocates a new rational with num, denom, and base.
func newRational(num, denom, base int64) *Rational {
	if denom == 0 {
		panic("Division by zero")
	}
	r := &Rational{big.NewInt(num), big.NewInt(denom), base}
	return r.normalize()
}

// Construct a new rational from a decimal.
//
// I, N, R: the value of the integer, non-repeating, repeating parts.
//
// n, r: the number of digits of the non-repeating and repeating parts.
func newRationalFromDec(I, N, R, n, r, base int64) *Rational {
	var num, denom int64

	if r == 0 {
		denom = intPow(base, int64(n))
		num = I*denom + int64(N)
	} else {
		x := (intPow(base, int64(r)) - 1)
		denom = intPow(base, int64(n)) * x
		num = I*denom + N*x + R
	}
	return newRational(num, denom, base)
}

// Sets n to the absolute value of n.
func (n *Rational) Abs() *Rational {
	if n.Sign() == -1 {
		n.num.Neg(n.num)
	}
	return n
}

// Compare this with another rational. Return 0 if equal, 1 if greater than,
// -1 if less than.
func (n *Rational) Cmp(other *Rational) int {
	matchDenom(n, other)
	return n.num.Cmp(other.num)
}

// Get the sign of this rational.
func (n *Rational) Sign() int {
	return n.num.Sign()
}

// Sets n to the result of n + other.
func (n *Rational) Add(other *Rational) *Rational {
	matchDenom(n, other)
	n.num.Add(n.num, other.num)
	return n.normalize()
}

// Split the rational into the quotient and remainder.
func (n *Rational) Divmod() (*big.Int, *big.Int) {
	q := new(big.Int).Div(n.num, n.denom)
	r := new(big.Int).Mod(n.num, n.denom)
	return q, r
}

// Sets n to the result of n - other.
func (n *Rational) Sub(other *Rational) *Rational {
	matchDenom(n, other)
	n.num.Sub(n.num, other.num)
	return n.normalize()
}

// Sets n to the result of n * other.
func (n *Rational) Mul(other *Rational) *Rational {
	n.num.Mul(n.num, other.num)
	n.denom.Mul(n.denom, other.denom)
	return n.normalize()
}

// Sets n to the result of n / other.
func (n *Rational) Div(other *Rational) (*Rational, error) {
	if other.num.Sign() == 0 {
		return nil, fmt.Errorf("Division by zero")
	}
	n.num.Mul(n.num, other.denom)
	n.denom.Mul(n.denom, other.num)
	return n.normalize(), nil
}

// Sets n to the result of -n.
func (n *Rational) Neg() *Rational {
	n.num.Neg(n.num)
	return n
}
