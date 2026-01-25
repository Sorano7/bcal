package mathb

import "fmt"

var Zero Rational = newRational(0, 1, 10)

// Represents a base rational in normalized form.
type Rational struct {
	num   int64
	denom int64
	base  int64
}

// Get the numerator of the rational.
func (n Rational) Num() int64 {
	return n.num
}

// Get the Denominator of the rational.
func (n Rational) Denom() int64 {
	return n.denom
}

// Get the (preferred) base of the rational.
func (n Rational) Base() int64 {
	return n.base
}

// Create a new rational with the provided base.
func (n Rational) WithBase(base int64) Rational {
	return newRational(n.num, n.denom, base)
}

// Constructs a new rational.
func newRational(num, denom, base int64) Rational {
	if denom == 0 {
		panic("Division by zero")
	}
	r := Rational{num, denom, base}
	return r.normalize()
}

// Construct a new rational from a decimal.
//
// I, N, R: the value of the integer, non-repeating, repeating parts.
//
// n, r: the number of digits of the non-repeating and repeating parts.
func newRationalFromDec(I, N, R, n, r, base int64) Rational {
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

// The absolute value of the rational.
func (n Rational) Abs() Rational {
	if n.Sign() == -1 {
		n.num = -n.num
	}
	return n
}

// Compare this with another rational. Return 0 if equal, 1 if greater than,
// -1 if less than.
func (n Rational) Cmp(other Rational) int {
	n, other = matchDenom(n, other)
	switch {
	case n.num == other.num:
		return 0
	case n.num > other.num:
		return 1
	default:
		return -1
	}
}

// Get the sign of this rational.
func (n Rational) Sign() int {
	switch {
	case n.num == 0:
		return 0
	case n.num > 0:
		return 1
	default:
		return -1
	}
}

// Add this to another rational.
func (n Rational) Add(other Rational) Rational {
	n, other = matchDenom(n, other)
	n.num += other.num
	return n.normalize()
}

// Split the rational into the quotient and remainder.
func (n Rational) Divmod() (int64, int64) {
	return n.num / n.denom, n.num % n.denom
}

func (n Rational) Sub(other Rational) Rational {
	n, other = matchDenom(n, other)
	n.num -= other.num
	return n.normalize()
}

func (n Rational) Mul(other Rational) Rational {
	n.num *= other.num
	n.denom *= other.denom
	return n.normalize()
}

func (n Rational) Div(other Rational) (Rational, error) {
	if other.num == 0 {
		return Rational{}, fmt.Errorf("Division by zero")
	}
	n.num *= other.denom
	n.denom *= other.num
	return n.normalize(), nil
}
