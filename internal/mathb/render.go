package mathb

import (
	"fmt"
	"math/big"
	"slices"
	"strings"
)

const (
	// The digit representations.
	Digits              = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	MaxBaseAlnum  int64 = int64(len(Digits))
	MaxFracDigits int   = 100
)

// Render the rational's internal representation for debugging.
func (n *Rational) Debug() string {
	return fmt.Sprintf("num: %v, denom: %v, base: %d", n.num, n.denom, n.base)
}

// Render the rational as p/q_b.
func (n *Rational) renderRational(alnum bool) string {
	num := renderIntInBase(n.num, n.base, alnum)
	denom := renderIntInBase(n.denom, n.base, alnum)
	return fmt.Sprintf("[%d](%s / %s)", n.base, num, denom)
}

// Render a rational to a string.
func (n *Rational) Render(useAlphaNum, inRational bool, maxFracDigits int) string {
	n.Normalize()
	if n.base > MaxBaseAlnum {
		useAlphaNum = false
	}
	if inRational {
		return n.renderRational(useAlphaNum)
	}
	return n.renderDecimal(useAlphaNum, maxFracDigits)
}

func (n *Rational) String() string {
	return n.Render(n.base <= MaxBaseAlnum, false, MaxFracDigits)
}

// Get the sign symbol for this rational.
func (n *Rational) signSymbol() string {
	sign := ""
	if n.Sign() == -1 {
		sign = "-"
	}
	return sign
}

// Convert an integer into a list of digits in base.
func intToList(n *big.Int, base int64) []int64 {
	var digits []int64
	b := new(big.Int)
	b.SetInt64(base)
	for n.Sign() > 0 {
		b.Mod(n, b)
		digits = append(digits, b.Int64())
		b.SetInt64(base)
		n.Div(n, b)
	}
	slices.Reverse(digits)
	return digits
}

// Render a int in the target base.
func renderIntInBase(n *big.Int, base int64, alnum bool) string {
	if base > MaxBaseAlnum {
		panic(fmt.Sprintf("Invalid base: %d", base))
	}
	if n.Sign() == 0 {
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

// Convert two lists of nonrepeating and repeating parts of the fraction to a string.
func fracToString(nonrep, rep []int64, base int64, alnum bool) string {
	var n, r string
	if alnum {
		n = foldToString(nonrep, base)
	} else {
		n = digitListString(nonrep)
	}
	if rep == nil {
		return fmt.Sprintf(".%s", n)
	}
	if alnum {
		r = foldToString(rep, base)
	} else {
		r = digitListString(rep)
	}
	return fmt.Sprintf(".%s(%s)", n, r)
}

// Render a terminating rational in its base digits.
func (n *Rational) renderDecimal(alnum bool, maxFracDigits int) string {
	var sb strings.Builder
	sb.WriteString(n.signSymbol())
	n = n.Clone()

	q, r := n.Abs().Divmod()
	sb.WriteString(renderIntInBase(q, n.base, alnum))
	if r.Sign() != 0 {
		nonrep, rep, ok := n.splitFrac(maxFracDigits)
		fmt.Fprint(&sb, fracToString(nonrep, rep, n.base, alnum))
		if !ok {
			sb.WriteString("...")
		}
	}
	return sb.String()
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
