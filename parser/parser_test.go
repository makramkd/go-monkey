package parser_test

import (
	"testing"

	"github.com/makramkd/go-monkey/ast"
	"github.com/makramkd/go-monkey/lexer"
	"github.com/makramkd/go-monkey/parser"
	"github.com/stretchr/testify/assert"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 3)
	assert.Empty(t, p.Errors())

	testCases := []struct {
		expectedIdentifier string
		expectedValue      string
	}{
		{"x", "5"},
		{"y", "10"},
		{"foobar", "838383"},
	}

	for i, testCase := range testCases {
		stmt := program.Statements[i]
		testLetStatement(t, stmt, testCase.expectedIdentifier)
	}
}

func TestLetStatementsError(t *testing.T) {
	input := `let a b`
	l := lexer.New(input)
	p := parser.New(l)
	p.ParseProgram()
	assert.Len(t, p.Errors(), 1)
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 10;
	return 828282;
	return add(15);`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 3)

	for _, stmt := range program.Statements {
		assert.IsType(t, &ast.ReturnStatement{}, stmt)
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) {
	assert.Equal(t, "let", stmt.TokenLiteral())
	assert.IsType(t, &ast.LetStatement{}, stmt)
	letStmt := stmt.(*ast.LetStatement)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	assert.Empty(t, p.Errors())
	assert.Len(t, program.Statements, 1)

	assert.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	assert.IsType(t, &ast.Identifier{}, stmt.Expression)

	ident := stmt.Expression.(*ast.Identifier)

	assert.Equal(t, "foobar", ident.Value)
	assert.Equal(t, "foobar", ident.TokenLiteral())
}
