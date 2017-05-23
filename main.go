package main

import (
	"fmt"
	"io/ioutil"
	"monkey/eval"
	"monkey/lexer"
	"monkey/parser"
	"monkey/repl"
	"os"
)

func runProgram(filename string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	f, err := ioutil.ReadFile(wd + "/" + filename)
	if err != nil {
		fmt.Println("monkey: ", err.Error())
		os.Exit(1)
	}
	l := lexer.New(string(f))
	p := parser.New(l, wd)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		fmt.Println(p.Errors()[0])
		os.Exit(1)
	}
	scope := eval.NewScope(nil)
	e := eval.Eval(program, scope)
	if e.Inspect() != "null" {
		fmt.Println(e.Inspect())
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Monkey programming language REPL\n")
		repl.Start(os.Stdout)
	} else {
		runProgram(args[0])
	}
}
