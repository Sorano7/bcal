package mathb

import (
	"fmt"
	"slices"
	"strings"
)

const (
	// The digit representations.
	Digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func (n Rational) RatString() string {
	return fmt.Sprintf("(%d/%d)_%d", n.num, n.denom, n.base)
}

func (n Rational) String() string {
	switch {
	case n.base >= int64(len(Digits)):
		return n.renderDecimalWithBase()
	case n.terminatesInBase():
		return n.renderTerminating()
	default:
		return n.renderRepeating()
	}
}

// Get the sign symbol for this rational.
func (n Rational) signSymbol() string {
	sign := ""
	if n.Sign() == -1 {
		sign = "-"
	}
	return sign
}

// Render a int in the target base.
func renderIntInBase(n, base int64) string {
	maxBase := len(Digits)
	if base > int64(maxBase) {
		panic(fmt.Sprintf("Invalid base: %d", base))
	}
	if n == 0 {
		return string(Digits[0])
	}
	var digits []int64
	for n > 0 {
		r := n % base
		n = n / base
		digits = append(digits, r)
	}
	slices.Reverse(digits)
	return foldToString(digits, base)
}

// Render a terminating rational in its base digits.
func (n Rational) renderTerminating() string {
	var sb strings.Builder
	sb.WriteString(n.signSymbol())

	q, r := n.Abs().Divmod()
	sb.WriteString(renderIntInBase(q, n.base))
	if r != 0 {
		nonrep, _ := n.splitFrac()
		fmt.Fprintf(&sb, ".%s", nonrep)
	}
	return sb.String()
}

// Render the rational as a base-10 decimal with base notation.
func (n Rational) renderDecimalWithBase() string {
	var sb strings.Builder
	sb.WriteString(n.signSymbol())

	q, r := n.Abs().Divmod()
	fmt.Fprint(&sb, q)
	if r != 0 {
		fmt.Fprintf(&sb, ".%d", r)
	}
	fmt.Fprintf(&sb, "(%d)", n.base)
	return sb.String()
}

// Render the rational as a repeating decimal in its base digits.
func (n Rational) renderRepeating() string {
	q, _ := n.Abs().Divmod()
	nonrep, rep := n.splitFrac()

	return fmt.Sprintf("%s.%s(%s)", renderIntInBase(q, n.base), nonrep, rep)
}

// Get the repeating part of this rational, converted to its base digits.
func (n Rational) splitFrac() (string, string) {
	_, r := n.Abs().Divmod()

	seen := make(map[int64]int)
	var digits []int64

	for r != 0 {
		seen[r] = len(digits)
		r *= n.base
		digits = append(digits, r/n.denom)
		r %= n.denom
		if _, ok := seen[r]; ok {
			break
		}
	}

	split := len(digits)
	if r != 0 {
		split = seen[r]
	}
	nonrep := digits[:split]
	rep := digits[split:]

	return foldToString(nonrep, n.base), foldToString(rep, n.base)
}

// Fold a slice of int to a string in the base.
func foldToString(slice []int64, base int64) string {
	var sb strings.Builder
	for _, n := range slice {
		sb.WriteString(valToChar(n, base))
	}
	return sb.String()
}

// Convert a value to its equivalent digit character.
func valToChar(v int64, base int64) string {
	if v >= base {
		panic(fmt.Sprintf("Invalid value %d for base %d", v, base))
	}
	return string(Digits[v])
}
