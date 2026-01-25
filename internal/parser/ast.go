package parser

import "fmt"

type Node interface {
	String() string
}

// Represents an AST statement.
type Statement interface {
	stmt()
	Node
}

// Represents an assignment.
type AssignStatement struct {
	Target Expression
	Value  Expression
}

func (a *AssignStatement) stmt() {}
func (a *AssignStatement) String() string {
	return fmt.Sprintf(`"assignment": {"target": {%s}, "value": {%s}}`, a.Target, a.Value)
}

// Represents a statement with only an expression.
type ExprStatement struct {
	Expr Expression
}

func (e *ExprStatement) stmt() {}
func (e *ExprStatement) String() string {
	return fmt.Sprintf(`"expr_statement": {"expr": {%s}}`, e.Expr)
}

// Represents an AST expression.
type Expression interface {
	expr()
	Node
}

// Represents an error node.
type ErrorNode struct {
	Msg string
}

func (e *ErrorNode) expr() {}
func (e *ErrorNode) stmt() {}
func (e *ErrorNode) String() string {
	return fmt.Sprintf(`"error": {"msg": "%s"}`, e.Msg)
}

// Create a new error node.
func newError(a any) *ErrorNode {
	return newErrorf("%s", a)
}

// Create a new error node with formatting.
func newErrorf(format string, a ...any) *ErrorNode {
	return &ErrorNode{Msg: fmt.Sprintf(format, a...)}
}

// Return whether the expression is an error.
func isError(n Node) bool {
	if n == nil {
		panic("Unexpected nil Node")
	}
	_, ok := n.(*ErrorNode)
	return ok
}

type DigitValue struct {
	Value int64
}

func (d *DigitValue) expr() {}
func (d *DigitValue) String() string {
	return fmt.Sprintf(`"digits": {"value": %d}`, d.Value)
}

type DigitString struct {
	Value string
}

func (d *DigitString) expr() {}
func (d *DigitString) String() string {
	return fmt.Sprintf(`"digits": {"value": "%s"}`, d.Value)
}

// Represents a number literal.
type NumberLiteral struct {
	Int    Expression
	Nonrep Expression
	Rep    Expression
}

func (d *NumberLiteral) expr() {}
func (d *NumberLiteral) String() string {
	return fmt.Sprintf(`"number": {"int": "%s", "nonrep": "%s", "rep": "%s"}`,
		d.Int, d.Nonrep, d.Rep)
}

type DigitList struct {
	Value []int64
}

func (d *DigitList) expr() {}
func (d *DigitList) String() string {
	return fmt.Sprintf(`"digits": {"value": %v}`, d.Value)
}

// Represents an identifier.
type Identifier struct {
	Name string
}

func (i *Identifier) expr() {}
func (i *Identifier) String() string {
	return fmt.Sprintf(`"ident": {"name": "%s"}`, i.Name)
}

// Represents an infix expression.
type InfixExpr struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpr) expr() {}
func (i *InfixExpr) String() string {
	return fmt.Sprintf(`"infix": {"left": {%s}, "op": "%s", "right": {%s}}`,
		i.Left, i.Operator, i.Right)
}

// Represents a prefix expression.
type PrefixExpr struct {
	Operator string
	Right    Expression
}

func (p *PrefixExpr) expr() {}
func (p *PrefixExpr) String() string {
	return fmt.Sprintf(`"prefix": {"op": "%s", "right": {%s}}`, p.Operator, p.Right)
}

// Represents a base annotation.
type BaseAnnotation struct {
	Base int64
	Expr Expression
}

func (b *BaseAnnotation) expr() {}
func (b *BaseAnnotation) String() string {
	return fmt.Sprintf(`"base_annotation": {"base": %d, "expr": {%s}}`, b.Base, b.Expr)
}

// Represents a output expression with base.
type OutputArguments struct {
	Args map[string]string
	Expr Expression
}

func (o *OutputArguments) expr() {}
func (o *OutputArguments) String() string {
	return fmt.Sprintf(`"output": {"args": %v, "expr": {%s}}`, o.Args, o.Expr)
}
