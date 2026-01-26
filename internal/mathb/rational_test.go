package mathb

import (
	"fmt"
	"math/big"
	"testing"
)

func bigIntEquals(a, b *big.Int) bool {
	return a.Cmp(b) == 0
}

func assertEquals(have, want *Rational, t *testing.T) {
	if have.Cmp(want) != 0 {
		t.Fatalf("have: %s, want: %s", have, want)
	}
}

func TestClone(t *testing.T) {
	x := newRational(big.NewInt(100), big.NewInt(1), 10)
	y := x.Clone()
	fmt.Printf("x(%p): { num: %p, denom: %p }\n", x, x.num, x.denom)
	fmt.Printf("y(%p): { num: %p, denom: %p }\n", y, y.num, y.denom)
}

func TestParseRender(t *testing.T) {
	str := "1A3B.2(56BA)"
	x, err := ParseString(str, 12)
	if err != nil {
		t.Fatal(err)
	}
	if x.Render(true, false) != str {
		t.Fatalf("Parsed value differs from render: %s", x.Render(true, false))
	}
}

func TestParseString(t *testing.T) {
	tmp := new(big.Int)
	x, err := ParseString("1A.25AB", 12)
	if err != nil {
		t.Fatal(err)
	}
	x.WithBase(10)

	if !bigIntEquals(x.num, tmp.SetInt64(460499)) && !bigIntEquals(x.denom, tmp.SetInt64(20736)) {
		t.Fatalf("Parse string failed: %s", x)
	}
}

func TestAdd(t *testing.T) {
	a, _ := ParseString("1B", 12)
	b, _ := ParseString("1", 10)
	want, _ := ParseString("20", 12)
	assertEquals(a.Add(b), want, t)
}
