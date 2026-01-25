package mathb

import (
	"fmt"
	"regexp"
	"strings"
)

// Parse a rational from parts.
func FromParts(intPart, nonrep, rep string, base int) (Rational, error) {
	if base >= len(Digits) {
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
		denom = intPow(int64(base), int64(n))
		num = I*denom + int64(N)
	} else {
		x := (intPow(int64(base), int64(r)) - 1)
		denom = intPow(int64(base), int64(n)) * x
		num = I*denom + N*x + R
	}

	return newRational(num, denom, int64(base)), nil
}

// Initialize a rational from a string.
func FromString(lit string, base int) (Rational, error) {
	if base >= len(Digits) {
		return Rational{}, fmt.Errorf("Base %d exceeds max representable base", base)
	}

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
func valueInBase(s string, base int) (int64, error) {
	if s == "" {
		return 0, nil
	}
	x := 0
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
func charToVal(c rune, base int) (int, error) {
	for i, digit := range Digits {
		if c != digit {
			continue
		}
		if i >= base {
			return 0, fmt.Errorf("Digit %c out of range for base %d", c, base)
		}
		return i, nil
	}
	return 0, fmt.Errorf("Invalid character: %c", c)
}
