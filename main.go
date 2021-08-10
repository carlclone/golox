package main

import (
	"bufio"
	"fmt"
	"os"
)

var hadError = false

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Fprint(os.Stderr, "usage:golox [script]\n")
		os.Exit(1)
	} else if len(args) == 2 {

	} else {
		runPrompt()
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		run(line)
		hadError = false
	}
}

func run(source string) {
	scanner := NewScanner(source)
	tokens, err := scanner.scan()
	if err != nil {
		fmt.Println(err)
		fmt.Println(tokens)
		hadError = true
		return
	}
	fmt.Println(tokens)
}
