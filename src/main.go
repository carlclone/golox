package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var hadError = false

type Local map[*Expr]int

func (l *Local) put(e Expr, depth int) {
	(*l)[&e] = depth
}
func (l *Local) get(e Expr) (int, bool) {
	v, ok := (*l)[&e]
	return v, ok
}

var locals = make(Local)

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Fprint(os.Stderr, "usage:golox [script]\n")
		os.Exit(1)
	} else if len(args) == 2 {
		runFile(args[1])
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

func runFile(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	run(string(data))
	if hadError {
		os.Exit(1)
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
	//for _, t := range tokens {
	//	fmt.Println("token ", t)
	//}

	p := NewParser(tokens)
	stmts, errs := p.parse()
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		hadError = true
		return
	}
	//for _, s := range stmts {
	//	fmt.Println(s)
	//	//printExprAST(s)
	//}

	resolver := NewResolver()
	resolver.resolve(stmts)

	globals := NewEnv(nil) // root env has no enclosure
	if err := interpret(stmts, globals); err != nil {
		fmt.Println(err)
		hadError = true
	}

}
