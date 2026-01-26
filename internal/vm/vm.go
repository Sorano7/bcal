package vm

import (
	"calculator/internal/parser"
)

// An VM for executing the program.
type VM struct {
	defaultBase int64
	store       map[string]Value
	lastVal     Value
}

// Create a new VM with the default base.
func newVM(defaultBase int64) *VM {
	vm := &VM{
		defaultBase: defaultBase,
		store:       make(map[string]Value),
	}
	return vm
}

// If the variable is defined.
func (v *VM) hasVar(name string) bool {
	_, ok := v.store[name]
	return ok
}

// Set the variable to a value.
func (v *VM) setVar(name string, value Value) {
	v.store[name] = value
}

// Get the variable value.
func (v *VM) getVar(name string) Value {
	if !v.hasVar(name) {
		return newErrorf("'%s' is not defined", name)
	}
	return v.store[name]
}

func (v *VM) run(src string) Value {
	program, err := parser.Parse(src)
	if err != nil {
		return newError(err)
	}
	return v.execute(program)
}
