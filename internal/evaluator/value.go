package evaluator

import (
	"calculator/internal/mathb"
	"fmt"
)

type ValueType int

const (
	_ ValueType = iota
	NumberValue
	ErrorValue
)

type Value interface {
	Type() ValueType
	String() string
}

type Number struct {
	Value mathb.Rational
}

func (n *Number) Type() ValueType { return NumberValue }
func (n *Number) String() string {
	return n.Value.String()
}

func newNumber(v mathb.Rational) *Number {
	return &Number{Value: v}
}

type Error struct {
	Msg string
}

func (e *Error) Type() ValueType { return ErrorValue }
func (e *Error) String() string {
	return e.Msg
}

func newError(a any) *Error {
	return newErrorf("%s", a)
}

func newErrorf(format string, a ...any) *Error {
	return &Error{Msg: fmt.Sprintf(format, a...)}
}

func isError(v Value) bool {
	return v.Type() == ErrorValue
}
