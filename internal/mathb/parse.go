package mathb

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

// Parse a rational from lists of digits.
func ParseDigitList(intPart, nonrep, rep []int64, base int64) (*Rational, error) {
	return newRationalFromDigits(intPart, nonrep, rep, base)
}

// Parse a rational from numerator and denominator with base respect.
func ParseRational(num, denom string, base int64) (*Rational, error) {
	n, err := stringValueInBase(num, base)
	if err != nil {
		return nil, err
	}
	d, err := stringValueInBase(denom, base)
	if err != nil {
		return nil, err
	}
	return newRational(n, d, base).Clone(), nil
}

// Parse a rational from parts.
func ParseParts(intPart, nonrep, rep string, base int64) (*Rational, error) {
	if base >= int64(len(Digits)) {
		return nil, fmt.Errorf("Base %d exceeds max representable base", base)
	}
	i, err := StringToList(intPart, base)
	if err != nil {
		return nil, err
	}
	n, err := StringToList(nonrep, base)
	if err != nil {
		return nil, err
	}
	r, err := StringToList(rep, base)
	if err != nil {
		return nil, err
	}
	return newRationalFromDigits(i, n, r, base)
}

// Initialize a rational from a string.
func ParseString(lit string, base int64) (*Rational, error) {
	parts := strings.Split(lit, ".")
	if len(parts) > 2 {
		return nil, fmt.Errorf("Invalid literal")
	}
	intPart := parts[0]
	nonrep := ""
	rep := ""
	if len(parts) == 2 {
		nonrep = parts[1]
		regex := regexp.MustCompile(`\([0-9A-Za-z]+\)`)
		match := regex.FindStringIndex(nonrep)
		if match != nil {
			rep = nonrep[match[0]+1 : match[1]-1]
			nonrep = nonrep[:match[0]]
		}
	}
	return ParseParts(intPart, nonrep, rep, base)
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
func stringValueInBase(s string, base int64) (*big.Int, error) {
	if s == "" {
		return nil, nil
	}
	digits, err := StringToList(s, base)
	if err != nil {
		return nil, err
	}
	return digitsToInt(digits, base)
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
