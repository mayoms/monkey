package main

import (
	"fmt"
	"monkey/repl"
	"os"
)

func main() {
	fmt.Println("Monkey programming language REPL\n")
	repl.Start(os.Stdout)
}
