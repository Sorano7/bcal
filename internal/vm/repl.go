package vm

import (
	"bufio"
	"calculator/internal/parser"
	"fmt"
	"os"
	"strings"
)

type Flow int

const (
	_ Flow = iota
	Continue
	Exit
)

// Start an interactive repl.
func StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	v := newVM(10)

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
			switch v.handleCommand(input[1:]) {
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
		fmt.Println(result)
	}
}

func (v *VM) handleCommand(cmd string) Flow {
	if cmd == "" {
		return Continue
	}
	parts := strings.Split(cmd, " ")
	cmd = parts[0]
	switch cmd {
	case "q", "quit":
		return Exit
	case "env":
		return v.printEnv()
	case "ast":
		if len(parts) == 1 {
			return Continue
		}
		return printAST(strings.Join(parts[1:], " "))
	default:
		fmt.Printf("Not a command: %s", cmd)
		return Exit
	}
}

func (v *VM) printEnv() Flow {
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
