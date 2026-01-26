package vm

import (
	"bufio"
	"calculator/internal/parser"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Flow int

const (
	_ Flow = iota
	Continue
	Exit
)

const MaxOutputChar = 100

type repl struct {
	maxOutput int
}

func newRepl() *repl {
	return &repl{MaxOutputChar}
}

// Start an interactive repl.
func StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	v := newVM(10)
	repl := newRepl()

MainLoop:
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if strings.TrimSpace(input) == "" {
			continue
		}

		if strings.HasPrefix(input, ":") {
			switch repl.handleCommand(input[1:], v) {
			case Continue:
				continue
			case Exit:
				break MainLoop
			}
		}

		result := v.run(input)
		if !isError(result) {
			v.lastVal = result
		}
		str := result.String()
		if repl.maxOutput > 0 && len(str) >= repl.maxOutput {
			str = str[:repl.maxOutput] + " (truncated)"
		}
		fmt.Println(str)
	}
}

func (r *repl) handleCommand(cmd string, v *VM) Flow {
	if cmd == "" {
		return Continue
	}
	parts := strings.Split(cmd, " ")
	cmd = parts[0]
	rest := ""
	if len(parts) > 1 {
		rest = strings.Join(parts[1:], " ")
	}
	switch cmd {
	case "q", "quit":
		return Exit
	case "env":
		return r.printEnv(v)
	case "ast":
		return printAST(rest)
	case "trun":
		n, err := strconv.ParseInt(rest, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		r.maxOutput = int(n)
		fmt.Printf("truncate=%d\n", r.maxOutput)
		return Continue
	default:
		fmt.Printf("Not a command: %s\n", cmd)
		return Continue
	}
}

func (r *repl) printEnv(v *VM) Flow {
	if len(v.store) == 0 {
		fmt.Println("empty")
		return Continue
	}
	for k, v := range v.store {
		fmt.Printf("@%s = %s\n", k, v)
	}
	return Continue
}

func printAST(src string) Flow {
	result, err := parser.Parse(src)
	if err != nil {
		fmt.Printf("%s\n", err)
		return Exit
	}
	fmt.Printf("%s\n", result)
	return Continue
}
