package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/makramkd/go-monkey/evaluator"
	"github.com/makramkd/go-monkey/lexer"
	"github.com/makramkd/go-monkey/object"
	"github.com/makramkd/go-monkey/parser"
	"github.com/makramkd/go-monkey/repl"
)

var stdlibModulesPath = flag.String("stdlib-modules-path", "", "The path to the Monkey standard library")
var executableFile = flag.String("e", "", "The Monkey script to execute and exit")

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	evaluator.SetStdlibPath(*stdlibModulesPath)

	if *executableFile != "" {
		executeMonkeyScript()
		return
	}

	fmt.Printf("Hello, %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

func executeMonkeyScript() {
	f, err := os.Open(*executableFile)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	l := lexer.New(string(b))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, e := range p.Errors() {
			fmt.Printf("parse error: %s", e.Error())
		}
		return
	}

	e := object.NewEnv()
	ret := evaluator.Eval(program, e)

	fmt.Printf("%s\n", ret.Inspect())
}
