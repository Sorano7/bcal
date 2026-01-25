package vm

import (
	"fmt"
	"slices"
	"strconv"
)

// A handler that transform the value based on the provided argument value.
type argHandler func(argv string, value Value) Value

// A map of handlers.
var handlers = map[string]argHandler{
	"base":  setOutputBase,
	"digit": setDigitsType,
	"num":   setNumericType,
	"prec":  nil,
}

// Get a handler for the argument name.
func getArgHandler(arg string) argHandler {
	if h, ok := handlers[arg]; ok {
		return h
	}
	return nil
}

// Assert the argument value is valid.
func validateArgv(have string, want ...string) error {
	if slices.Contains(want, have) {
		return nil
	}
	return fmt.Errorf("Invalid argument value: %s", have)
}

// Set the output base.
func setOutputBase(argv string, value Value) Value {
	base, err := strconv.ParseInt(argv, 10, 64)
	if err != nil {
		return newError(err)
	}
	switch v := value.(type) {
	case *Number:
		v.Value = v.Value.WithBase(base)
		return v
	default:
		return v
	}
}

// Set the rendering type of the digits.
func setDigitsType(argv string, value Value) Value {
	if err := validateArgv(argv, "alnum", "list"); err != nil {
		return newError(err)
	}
	switch v := value.(type) {
	case *Number:
		v.UseAlnum = argv == "alnum"
		return v
	default:
		return v
	}
}

// Set the numeric type of the number.
func setNumericType(argv string, value Value) Value {
	if err := validateArgv(argv, "rational", "decimal"); err != nil {
		return newError(err)
	}
	switch v := value.(type) {
	case *Number:
		v.inRational = argv == "rational"
		return v
	default:
		return v
	}
}
