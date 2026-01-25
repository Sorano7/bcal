package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Start an interactive repl.
func StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if strings.TrimSpace(input) == "" {
			continue
		}
		if input == ":q" {
			break
		}

		e := newVM(10)
		result, err := e.Evaluate(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		fmt.Println(result)
	}
}
