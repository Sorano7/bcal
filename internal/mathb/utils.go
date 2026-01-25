package mathb

import "math/big"

// Computes the GCD of two integers.
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Computes the LCM of the provided intergers.
func lcm(a, b *big.Int) *big.Int {
	g := new(big.Int).GCD(nil, nil, a, b)
	t := new(big.Int).Div(a, g)
	t.Abs(t)
	t.Mul(t, b)
	t.Abs(a)
	return t
}

// Mutates and normalize the rational.
func (n *Rational) normalize() *Rational {
	g := new(big.Int).GCD(nil, nil, n.num, n.denom)
	n.num.Div(n.num, g)
	n.denom.Div(n.denom, g)
	if n.denom.Sign() < 0 {
		n.num.Neg(n.num)
		n.denom.Neg(n.denom)
	}
	return n
}

// Mutates to match the two rationals' denominator and return both rationals.
func matchDenom(a, b *Rational) {
	d := lcm(a.denom, b.denom)
	tmp := new(big.Int)
	a.num.Mul(a.num, tmp.Div(d, a.denom))
	b.num.Mul(b.num, tmp.Div(d, b.denom))
	a.denom = d
	tmp.Set(d)
	b.denom = tmp
}

// Integer power.
func intPow(a, b int64) int64 {
	var x int64 = 1
	for b > 0 {
		if b%2 == 1 {
			x *= a
		}
		a *= a
		b /= 2
	}
	return x
}

// Compute the non-repeating and repeating part of this rational.
func (n *Rational) splitFrac() ([]int64, []int64) {
	tmp := *n
	denom := tmp.denom

	rem := tmp.Abs().num.Mod(tmp.num, denom)
	if rem.Sign() == 0 {
		return nil, nil
	}

	b := big.NewInt(n.base)
	seen := make(map[string]int, 64)
	digits := make([]int64, 0, 64)

	t := new(big.Int)
	for {
		if rem.Sign() == 0 {
			return digits, nil
		}

		key := string(rem.Bytes())
		if split, ok := seen[key]; ok {
			return digits[:split], digits[split:]
		}
		seen[key] = len(digits)

		t.Mul(rem, b)
		t.QuoRem(t, denom, rem)
		digits = append(digits, t.Int64())
	}
}
