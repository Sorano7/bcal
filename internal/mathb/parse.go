package mathb

import (
	"fmt"
	"regexp"
	"strings"
)

// Parse a rational from lists of digits.
func ParseDigitList(intPart, nonrep, rep []int64, base int64) (Rational, error) {
	I, err := basebAcc(intPart, base)
	if err != nil {
		return Rational{}, err
	}
	N, err := basebAcc(nonrep, base)
	if err != nil {
		return Rational{}, err
	}
	R, err := basebAcc(rep, base)
	if err != nil {
		return Rational{}, err
	}

	n, r := int64(len(nonrep)), int64(len(rep))

	return newRationalFromDec(I, N, R, n, r, base), nil
}

// Parse a rational from numerator and denominator with base respect.
func ParseRational(num, denom string, base int64) (Rational, error) {
	n, err := valueInBase(num, base)
	if err != nil {
		return Rational{}, err
	}
	d, err := valueInBase(denom, base)
	if err != nil {
		return Rational{}, err
	}
	return newRational(n, d, base), nil
}

// Parse a rational from parts.
func ParseParts(intPart, nonrep, rep string, base int64) (Rational, error) {
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

	n := int64(len(nonrep))
	r := int64(len(rep))

	return newRationalFromDec(I, N, R, n, r, base), nil
}

// Initialize a rational from a string.
func ParseString(lit string, base int64) (Rational, error) {
	parts := strings.Split(lit, ".")
	if len(parts) > 2 {
		return Rational{}, fmt.Errorf("Invalid literal")
	}
	intPart := parts[0]
	nonrep := ""
	rep := ""
	if len(parts) == 2 {
		nonrep = parts[1]
		regex := regexp.MustCompile(`\([0-9A-Za-z]+\)`)
		match := regex.FindStringIndex(parts[1])
		if match != nil {
			rep = nonrep[match[0]:match[1]]
			nonrep = nonrep[:match[0]]
		}
	}
	return ParseParts(intPart, nonrep, rep, base)
}

// Base-b accummulation.
func basebAcc(digits []int64, base int64) (int64, error) {
	if digits == nil {
		return 0, nil
	}
	var v int64 = 0
	for _, d := range digits {
		if d >= base {
			return 0, fmt.Errorf("Digit %d out of range for base %d", d, base)
		}
		v = v*base + d
	}
	return v, nil
}

func StringToList(digits string, base int64) ([]int64, error) {
	if digits == "" {
		return nil, nil
	}
	var out []int64
	for _, c := range digits {
		val, err := charToVal(c, base)
		if err != nil {
			return nil, err
		}
		out = append(out, val)
	}
	return out, nil
}

// Parse the value of a string in the given base.
func valueInBase(s string, base int64) (int64, error) {
	if s == "" {
		return 0, nil
	}
	digits, err := StringToList(s, base)
	if err != nil {
		return 0, err
	}
	return basebAcc(digits, base)
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
