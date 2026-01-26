package mathb

import (
	"fmt"
	"math/big"
)

const (
	MaxRepeatDetect int = 1e6
)

// Represents a base rational in.Normalized form.
type Rational struct {
	num   *big.Int
	denom *big.Int
	base  int64
}

// Creates a deep copy.
func (n *Rational) Clone() *Rational {
	num := new(big.Int).Set(n.num)
	denom := new(big.Int).Set(n.denom)
	return &Rational{num, denom, n.base}
}

// Sets n's base to base.
func (n *Rational) WithBase(base int64) *Rational {
	n.base = base
	return n
}

// Allocates a new rational with num, denom, and base.
func newRational(num, denom *big.Int, base int64) *Rational {
	if denom.Sign() == 0 {
		panic("Division by zero")
	}
	n := &Rational{num, denom, base}
	return n.Normalize()
}

// Construct a new rational from a decimal. Input will be mutated.
func newRationalFromDigits(intPart, nonrep, rep []int64, base int64) (*Rational, error) {
	I, err := digitsToInt(intPart, base)
	if err != nil {
		return nil, err
	}
	N, err := digitsToInt(nonrep, base)
	if err != nil {
		return nil, err
	}
	R, err := digitsToInt(rep, base)
	if err != nil {
		return nil, err
	}

	var num, denom *big.Int

	m, k := len(nonrep), len(rep)
	Bm := intPow(base, int64(m))

	if k == 0 {
		num = new(big.Int).Mul(I, Bm)
		num.Add(num, N)
		denom = new(big.Int).Set(Bm)
	} else {
		tmp := intPow(base, int64(k))
		tmp.Sub(tmp, big.NewInt(1))
		denom = new(big.Int).Mul(Bm, tmp)
		num = new(big.Int).Mul(I, denom)

		tmp.Mul(N, tmp)
		num.Add(num, tmp)
		num.Add(num, R)
	}
	return newRational(num, denom, base).Clone(), nil
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
	return n
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
	return n
}

// Sets n to the result of n * other.
func (n *Rational) Mul(other *Rational) *Rational {
	n.num.Mul(n.num, other.num)
	n.denom.Mul(n.denom, other.denom)
	return n
}

// Sets n to the result of n / other.
func (n *Rational) Div(other *Rational) (*Rational, error) {
	if other.num.Sign() == 0 {
		return nil, fmt.Errorf("Division by zero")
	}
	n.num.Mul(n.num, other.denom)
	n.denom.Mul(n.denom, other.num)
	return n, nil
}

// Sets n to the result of -n.
func (n *Rational) Neg() *Rational {
	n.num.Neg(n.num)
	return n
}
