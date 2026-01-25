package mathb_test

import (
	"calculator/internal/mathb"
	"fmt"
	"testing"
)

func assertEqual(value, want mathb.Rational, t *testing.T) {
	if value.Cmp(want) != 0 {
		t.Fatalf("want: %s (%s), got: %s (%s)",
			want, want.Render(false, true), value, value.Render(false, true))
	}
}

func printExpr(a, b, result mathb.Rational) {
	fmt.Printf("%s <op> %s = %s\n", a.Render(false, true), b.Render(false, true), result.Render(false, true))
}

func TestParsingDigits(t *testing.T) {
	x, _ := mathb.ParseString("ABC", 16)
	y, _ := mathb.ParseDigitList([]int64{10, 11, 12}, nil, nil, 16)
	assertEqual(x, y, t)
}

func TestParsingString(t *testing.T) {
	if _, err := mathb.ParseString("100", 10); err != nil {
		t.Fatal(err)
	}
	if _, err := mathb.ParseString("1234", 2); err == nil {
		t.Fatal(err)
	}
	if _, err := mathb.ParseString("&&&", 60); err == nil {
		t.Fatal(err)
	}
	if _, err := mathb.ParseString("ABC", 100); err == nil {
		t.Fatal(err)
	}
	str := "1AB24.AB5345"
	x, _ := mathb.ParseString(str, 12)
	if x.String() != str {
		t.Fatalf("want: %s, got: %s", str, x)
	}
}

func TestAddInt(t *testing.T) {
	x, _ := mathb.ParseString("1A", 12)
	y, _ := mathb.ParseString("1", 12)
	want, _ := mathb.ParseString("1B", 12)
	assertEqual(x.Add(y), want, t)
}

func TestAddFrac(t *testing.T) {
	x, _ := mathb.ParseString("1.A1", 12)
	y, _ := mathb.ParseString("2.3B", 12)
	want, _ := mathb.ParseString("4.2", 12)
	assertEqual(x.Add(y), want, t)
}
