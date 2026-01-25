package vm

import (
	"calculator/internal/mathb"
	"calculator/internal/parser"
)

// Evaluate an expression.
func (v *VM) evalExpr(expr parser.Expression, base int64) Value {
	switch e := expr.(type) {
	case *parser.BaseAnnotation:
		return v.evalIOBase(e.Expr, e.Base, e.Base)
	case *parser.OutputBase:
		return v.evalIOBase(e.Expr, base, e.Base)
	case *parser.NumberLiteral:
		return v.evalNumber(e, base)
	case *parser.PrefixExpr:
		return v.evalPrefix(e, base)
	case *parser.InfixExpr:
		return v.evalInfix(e, base)
	case *parser.Identifier:
		return v.evalIdent(e)
	default:
		return newErrorf("Expression not supported: %T", expr)
	}
}

// Evaluate with input and output base.
func (v *VM) evalIOBase(expr parser.Expression, ibase, obase int64) Value {
	value := v.evalExpr(expr, ibase)
	if isError(value) {
		return value
	}
	switch val := value.(type) {
	case *Number:
		return newNumber(val.Value.WithBase(obase))
	default:
		return newErrorf("Invalid expr: %s", val)
	}
}

// Evaluate a number expression.
func (ev *VM) evalNumber(e *parser.NumberLiteral, base int64) Value {
	v, err := mathb.FromParts(e.Int, e.Nonrep, e.Rep, base)
	if err != nil {
		return newError(err)
	}
	return newNumber(v)
}

// Evaluate a prefix expression.
func (v *VM) evalPrefix(e *parser.PrefixExpr, base int64) Value {
	right := v.evalExpr(e.Right, base)
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
func (v *VM) evalInfix(e *parser.InfixExpr, base int64) Value {
	left := v.evalExpr(e.Left, base)
	if isError(left) {
		return left
	}
	right := v.evalExpr(e.Right, base)
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

func (v *VM) evalIdent(e *parser.Identifier) Value {
	if e.Name != "" {
		return v.getVar(e.Name)
	}
	if v.lastVal != nil {
		return v.lastVal
	}
	return newError("Undefined")
}
