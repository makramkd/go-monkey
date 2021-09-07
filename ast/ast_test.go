package ast_test

import (
	"testing"

	"github.com/makramkd/go-monkey/ast"
	"github.com/makramkd/go-monkey/token"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.New(token.LET, "let"),
				Name: &ast.Identifier{
					Token: token.New(token.IDENT, "myVar"),
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.New(token.IDENT, "anotherVar"),
					Value: "anotherVar",
				},
			},
		},
	}

	assert.Equal(t, "let myVar = anotherVar;", program.String())
}
