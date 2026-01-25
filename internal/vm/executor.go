package vm

import "calculator/internal/parser"

// Execute a statement.
func (v *VM) execute(stmt parser.Statement) Value {
	switch s := stmt.(type) {
	case *parser.AssignStatement:
		return v.handleAssign(s)
	case *parser.ExprStatement:
		return v.evalExpr(s.Expr, v.defaultBase)
	default:
		return newErrorf("Statement not supported: %T", stmt)
	}
}

// Handle assignment.
func (v *VM) handleAssign(e *parser.AssignStatement) Value {
	switch t := e.Target.(type) {
	case *parser.Identifier:
		value := v.evalExpr(e.Value, v.defaultBase)
		if !isError(value) {
			v.setVar(t.Name, value)
		}
		return value
	default:
		return newErrorf("Invalid assignment target: %T", t)
	}
}
