package mathb

import (
	"fmt"
	"math/big"
)

// Computes the LCM of the provided intergers.
func lcm(a, b *big.Int) *big.Int {
	g := new(big.Int).GCD(nil, nil, a, b)
	t := new(big.Int).Div(a, g)
	t.Abs(t)
	t.Mul(t, b)
	t.Abs(a)
	return t
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
func intPow(a, b int64) *big.Int {
	return new(big.Int).Exp(big.NewInt(a), big.NewInt(b), nil)
}

// Compute the non-repeating and repeating part of this rational.
func (n *Rational) splitFrac(maxDigit int) ([]int64, []int64, bool) {
	n = n.Clone()
	denom := n.denom

	rem := n.Abs().num.Mod(n.num, denom)
	if rem.Sign() == 0 {
		return nil, nil, true
	}

	b := big.NewInt(n.base)
	seen := make(map[string]int, 64)
	digits := make([]int64, 0, 64)

	t := new(big.Int)
	for i := 0; i < MaxRepeatDetect; i++ {
		if rem.Sign() == 0 {
			return digits, nil, true
		}
		if maxDigit > 0 && len(digits) >= maxDigit {
			break
		}

		key := string(rem.Bytes())
		if split, ok := seen[key]; ok {
			return digits[:split], digits[split:], true
		}
		seen[key] = len(digits)

		t.Mul(rem, b)
		t.QuoRem(t, denom, rem)
		digits = append(digits, t.Int64())
	}
	return digits, nil, false
}

// Get the value of the list of digits in the provided base.
func digitsToInt(digits []int64, base int64) (*big.Int, error) {
	b := big.NewInt(base)
	v, tmp := new(big.Int), new(big.Int)

	for _, d := range digits {
		if d < 0 || d >= base {
			return nil, fmt.Errorf("Digit %d out of range for base %d", d, base)
		}
		v.Mul(v, b)
		v.Add(v, tmp.SetInt64(d))
	}
	return v, nil
}
