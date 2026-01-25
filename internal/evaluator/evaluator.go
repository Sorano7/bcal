package evaluator

import (
	"calculator/internal/mathb"
	"calculator/internal/parser"
	"fmt"
)

// An VM for executing the program.
type VM struct {
	defaultBase int
}

// Create a new VM with the default base.
func newVM(defaultBase int) *VM {
	return &VM{defaultBase}
}

// Evaluate the source with base.
func (v *VM) Evaluate(src string) (Value, error) {
	expr, err := parser.Parse(src)
	if err != nil {
		return nil, err
	}
	result := v.evalExpr(expr)
	if isError(result) {
		return nil, fmt.Errorf("%s", result.String())
	}
	return result, nil
}

// Evaluate an expression.
func (v *VM) evalExpr(expr parser.Expression) Value {
	switch e := expr.(type) {
	case *parser.NumberLiteral:
		return v.evalNumber(e)
	case *parser.PrefixExpr:
		return v.evalPrefix(e)
	case *parser.InfixExpr:
		return v.evalInfix(e)
	default:
		return newErrorf("Expression not supported: %T", expr)
	}
}

// Evaluate a number expression.
func (ev *VM) evalNumber(e *parser.NumberLiteral) Value {
	if e.Base == 0 {
		e.Base = ev.defaultBase
	}
	v, err := mathb.FromParts(e.Int, e.Nonrep, e.Rep, e.Base)
	if err != nil {
		return newError(err)
	}
	return newNumber(v)
}

// Evaluate a prefix expression.
func (v *VM) evalPrefix(e *parser.PrefixExpr) Value {
	right := v.evalExpr(e.Right)
	if isError(right) {
		return right
	}
	switch {
	case right.Type() == NumberValue:
		return v.evalNumberPrefix(e.Operator, right.(*Number))
	default:
		return newErrorf("Invalid operation: %s%s", e.Operator, right)
	}
}

// Evaluate a number prefix.
func (v *VM) evalNumberPrefix(op string, n *Number) Value {
	switch op {
	default:
		return newErrorf("Invalid prefix operation: %s<number>", op)
	}
}

// Evaluate an infix expression.
func (v *VM) evalInfix(e *parser.InfixExpr) Value {
	left := v.evalExpr(e.Left)
	if isError(left) {
		return left
	}
	right := v.evalExpr(e.Right)
	if isError(right) {
		return right
	}
	switch {
	case left.Type() == NumberValue && right.Type() == NumberValue:
		return v.evalNumberInfix(e.Operator, left.(*Number), right.(*Number))
	default:
		return newErrorf("Invalid operation: %s %s %s", left, e.Operator, right)
	}
}

// Evaluate a number infix.
func (v *VM) evalNumberInfix(op string, left, right *Number) Value {
	l, r := left.Value, right.Value
	switch op {
	case "+":
		return newNumber(l.Add(r))
	case "-":
		return newNumber(l.Sub(r))
	case "*":
		return newNumber(l.Mul(r))
	case "/":
		res, err := l.Div(r)
		if err != nil {
			return newError(err)
		}
		return newNumber(res)
	default:
		return newErrorf("Invalid prefix operation: <number> %s <number>", op)
	}
}
