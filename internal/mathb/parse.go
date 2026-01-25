package mathb

import (
	"fmt"
	"regexp"
	"strings"
)

// Parse a rational from parts.
func FromParts(intPart, nonrep, rep string, base int64) (Rational, error) {
	if base >= int64(len(Digits)) {
		return Rational{}, fmt.Errorf("Base %d exceeds max representable base", base)
	}
	I, err := valueInBase(intPart, base)
	if err != nil {
		return Rational{}, err
	}
	N, err := valueInBase(nonrep, base)
	if err != nil {
		return Rational{}, err
	}
	R, err := valueInBase(rep, base)
	if err != nil {
		return Rational{}, err
	}

	n := len(nonrep)
	r := len(rep)

	var num, denom int64

	if r == 0 {
		denom = intPow(base, int64(n))
		num = I*denom + int64(N)
	} else {
		x := (intPow(base, int64(r)) - 1)
		denom = intPow(base, int64(n)) * x
		num = I*denom + N*x + R
	}

	return newRational(num, denom, base), nil
}

// Initialize a rational from a string.
func FromString(lit string, base int64) (Rational, error) {
	parts := strings.Split(lit, ".")
	if len(parts) > 2 {
		return Rational{}, fmt.Errorf("Invalid literal")
	}
	intPart := parts[0]
	nonrep := ""
	rep := ""
	if len(parts) == 2 {
		frac := parts[1]
		regex := regexp.MustCompile(`\[[0-9A-Za-z]+\]`)
		match := regex.FindStringIndex(parts[1])
		if match != nil {
			rep = frac[match[0]:match[1]]
		}
		nonrep = frac[:match[0]]
	}
	return FromParts(intPart, nonrep, rep, base)
}

// Parse the value of a string in the given base.
func valueInBase(s string, base int64) (int64, error) {
	if s == "" {
		return 0, nil
	}
	var x int64 = 0
	for _, c := range s {
		val, err := charToVal(c, base)
		if err != nil {
			return 0, err
		}
		x = x*base + val
	}
	return int64(x), nil
}

// Convert a character to digit value.
func charToVal(c rune, base int64) (int64, error) {
	for i, digit := range Digits {
		if c != digit {
			continue
		}
		if int64(i) >= base {
			return 0, fmt.Errorf("Digit %c (%d) out of range for base %d", c, i, base)
		}
		return int64(i), nil
	}
	return 0, fmt.Errorf("Invalid character: %c", c)
}
