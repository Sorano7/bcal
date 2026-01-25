package mathb_test

import (
	"calculator/internal/mathb"
	"fmt"
	"testing"
)

func assertEqual(value, want mathb.Rational, t *testing.T) {
	if value.Cmp(want) != 0 {
		t.Fatalf("want: %s (%s), got: %s (%s)",
			want, want.RatString(), value, value.RatString())
	}
}

func printExpr(a, b, result mathb.Rational) {
	fmt.Printf("%s <op> %s = %s\n", a.RatString(), b.RatString(), result.RatString())
}

func TestParsing(t *testing.T) {
	if _, err := mathb.FromString("100", 10); err != nil {
		t.Fatal(err)
	}
	if _, err := mathb.FromString("1234", 2); err == nil {
		t.Fatal(err)
	}
	if _, err := mathb.FromString("&&&", 60); err == nil {
		t.Fatal(err)
	}
	if _, err := mathb.FromString("ABC", 100); err == nil {
		t.Fatal(err)
	}
	str := "1AB24.AB5345"
	x, _ := mathb.FromString(str, 12)
	if x.String() != str {
		t.Fatalf("want: %s, got: %s", str, x)
	}
}

func TestAddInt(t *testing.T) {
	x, _ := mathb.FromString("1A", 12)
	y, _ := mathb.FromString("1", 12)
	want, _ := mathb.FromString("1B", 12)
	assertEqual(x.Add(y), want, t)
}

func TestAddFrac(t *testing.T) {
	x, _ := mathb.FromString("1.A1", 12)
	y, _ := mathb.FromString("2.3B", 12)
	want, _ := mathb.FromString("4.2", 12)
	assertEqual(x.Add(y), want, t)
}
