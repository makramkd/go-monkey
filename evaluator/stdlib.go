package evaluator

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/makramkd/go-monkey/ast"
	"github.com/makramkd/go-monkey/lexer"
	"github.com/makramkd/go-monkey/parser"
)

var stdlibPath = "stdlib"

func SetStdlibPath(path string) {
	stdlibPath = path
}

func loadStdModule(name string) (*ast.Program, error) {
	f, err := os.Open(fmt.Sprintf("%s/%s.monkey", stdlibPath, name))
	if err != nil {
		return nil, fmt.Errorf("standard module does not exist: %s. Was the Monkey stdlib path specified correctly?", name)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("internal error: failed to read standard module '%s': %v", name, err)
	}

	// NOTE: we expect standard modules to be free of errors :)
	l := lexer.New(string(b))
	p := parser.New(l)
	return p.ParseProgram(), nil
}
