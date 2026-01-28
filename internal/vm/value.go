package vm

import (
	"calculator/internal/mathb"
	"fmt"
)

// Represents the value type.
type ValueType int

const (
	_ ValueType = iota
	NumberValue
	ErrorValue
	VoidValue
)

// Represents a runtime value.
type Value interface {
	Type() ValueType
	String() string
}

// Represents a number.
type Number struct {
	Value      *mathb.Rational
	UseAlnum   bool
	inRational bool
	Prec       int
}

func (n *Number) Type() ValueType { return NumberValue }
func (n *Number) String() string {
	return n.Value.Render(n.UseAlnum, n.inRational, n.Prec)
}

// Create a new number with a clone of n.
func newNumber(n *mathb.Rational) *Number {
	return &Number{
		Value:      n.Clone(),
		UseAlnum:   true,
		inRational: false,
		Prec:       20,
	}
}

// Represents a runtime error.
type Error struct {
	Msg string
}

func (e *Error) Type() ValueType { return ErrorValue }
func (e *Error) String() string {
	return e.Msg
}

// Create a new error with the default formatting.
func newError(a any) *Error {
	return newErrorf("%s", a)
}

// Create a new error with formatting.
func newErrorf(format string, a ...any) *Error {
	return &Error{Msg: fmt.Sprintf(format, a...)}
}

// Returns whether v is a runtime error.
func isError(v Value) bool {
	return v.Type() == ErrorValue
}

// Represents a void value.
type Void struct{}

func (v *Void) Type() ValueType { return VoidValue }
func (v *Void) String() string  { return "" }
