package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/makramkd/go-monkey/evaluator"
	"github.com/makramkd/go-monkey/repl"
)

var stdlibModulesPath = flag.String("stdlib-modules-path", "", "The path to the Monkey standard library")

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	evaluator.SetStdlibPath(*stdlibModulesPath)

	fmt.Printf("Hello, %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
