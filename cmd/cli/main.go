package main

import (
	"calculator/internal/vm"
	"flag"
	"fmt"
	"strings"
)

func main() {
	showHelpPtr := flag.Bool("help", false, "Display manual")

	flag.Parse()
	args := flag.Args()
	if *showHelpPtr {
		printHelp()
	} else if len(args) > 0 {
		result := vm.Run(strings.Join(args, " "))
		fmt.Println(result)
	} else {
		vm.StartREPL()
	}
}

func printHelp() {
	msg := `
Usage:
	bcal                    - Start REPL
	bcal --help             - Show this page
	bcal <expr>             - Evaluate the expression and print the result

Syntax:
	Base Syntax             - [<base>]<value>.(<value>) #<args>, ...
	Alphanumerics           - 0-9...A-Z...a-z
	List                    - {..., ..., ...}

Variables:
	@<ident> = ...          - Bind a variable
	@                       - The last result

Output:
	base=<value>            - The output base 
	digit=<alnum|list>      - The digit type (alphanumerics or list of digits)
	num=<decimal|rational>  - The number type (decimal or rational)
	prec=<value>            - Precision; 0 to display as repeat.

REPL:
	:q, :quit               - Quit the REPL
	:env                    - Show the current environment store
	:trun <value>           - Sets the maximum output character; 0 to disable
	`
	fmt.Print(msg)

}
