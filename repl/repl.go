package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/makramkd/go-monkey/evaluator"
	"github.com/makramkd/go-monkey/lexer"
	"github.com/makramkd/go-monkey/object"
	"github.com/makramkd/go-monkey/parser"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnv()

	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		// NOTE: evaluated == nil doesn't mean that there's an error. It just means we've executed
		// a statement that has no output.
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		io.WriteString(out, "\t"+err.Error()+"\n")
	}
}
