package mathb

import (
	"fmt"
	"slices"
	"strings"
)

const (
	// The digit representations.
	Digits             = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	MaxBaseAlnum int64 = int64(len(Digits))
)

// Render the rational as p/q_b.
func (n Rational) renderRational(alnum bool) string {
	num := renderIntInBase(n.num, n.base, alnum)
	denom := renderIntInBase(n.denom, n.base, alnum)
	return fmt.Sprintf("[%d](%s / %s)", n.base, num, denom)
}

// Render a rational to a string.
func (n Rational) Render(useAlphaNum, inRational bool) string {
	if n.base > MaxBaseAlnum {
		useAlphaNum = false
	}
	if inRational {
		return n.renderRational(useAlphaNum)
	}
	if n.terminatesInBase() {
		return n.renderTerminating(useAlphaNum)
	}
	return n.renderRepeating(useAlphaNum)
}

func (n Rational) String() string {
	return n.Render(n.base <= MaxBaseAlnum, false)
}

// Get the sign symbol for this rational.
func (n Rational) signSymbol() string {
	sign := ""
	if n.Sign() == -1 {
		sign = "-"
	}
	return sign
}

// Convert an integer into a list of digits in base.
func intToList(n, base int64) []int64 {
	var digits []int64
	for n > 0 {
		r := n % base
		n = n / base
		digits = append(digits, r)
	}
	slices.Reverse(digits)
	return digits
}

// Render a int in the target base.
func renderIntInBase(n, base int64, alnum bool) string {
	if base > MaxBaseAlnum {
		panic(fmt.Sprintf("Invalid base: %d", base))
	}
	if n == 0 {
		return string(Digits[0])
	}
	digits := intToList(n, base)
	if alnum {
		return foldToString(digits, base)
	}
	return digitListString(digits)
}

// Convert a digit list to a string.
func digitListString(digits []int64) string {
	var strs []string
	for _, d := range digits {
		strs = append(strs, fmt.Sprint(d))
	}
	return fmt.Sprintf("{%s}", strings.Join(strs, ", "))
}

// Render a terminating rational in its base digits.
func (n Rational) renderTerminating(alnum bool) string {
	var sb strings.Builder
	sb.WriteString(n.signSymbol())

	q, r := n.Abs().Divmod()
	sb.WriteString(renderIntInBase(q, n.base, alnum))
	if r != 0 {
		nonrep, _ := n.splitFrac(alnum)
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
func (n Rational) renderRepeating(alnum bool) string {
	q, _ := n.Abs().Divmod()
	nonrep, rep := n.splitFrac(alnum)

	return fmt.Sprintf("%s.%s(%s)", renderIntInBase(q, n.base, alnum), nonrep, rep)
}

// Get the repeating part of this rational, converted to its base digits or list.
func (n Rational) splitFrac(alnum bool) (string, string) {
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

	if alnum {
		return foldToString(nonrep, n.base), foldToString(rep, n.base)
	}
	return digitListString(nonrep), digitListString(rep)
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
