package mathb

// Computes the GCD of two integers.
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Computes the LCM of the provided intergers.
func lcm(a, b int64, n ...int64) int64 {
	res := (a / gcd(a, b)) * b
	for i := range n {
		res = lcm(res, n[i])
	}
	return res
}

// Normalize the rational.
func (n Rational) normalize() Rational {
	d := gcd(n.num, n.denom)
	n = Rational{n.num / d, n.denom / d, n.base}
	if n.denom < 0 {
		n = Rational{-n.num, -n.denom, n.base}
	}
	return n
}

// Match the two rationals' denominator and return both rationals.
func matchDenom(a, b Rational) (Rational, Rational) {
	d := lcm(a.denom, b.denom)
	a = Rational{a.num * (d / a.denom), d, a.base}
	b = Rational{b.num * (d / b.denom), d, b.base}
	return a, b
}

// Return whether this rational terminates in its prefered base.
func (n Rational) terminatesInBase() bool {
	n = n.Abs()

	g := gcd(n.num, n.denom)
	n.denom /= g

	for {
		g := gcd(n.denom, n.base)
		if g == 1 {
			break
		}
		n.denom /= g
	}
	return n.denom == 1
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
