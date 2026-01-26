package main

import (
	"calculator/internal/vm"
	"flag"
	"fmt"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		result := vm.Run(strings.Join(args, " "))
		fmt.Println(result)
	} else {
		vm.StartREPL()
	}
}
