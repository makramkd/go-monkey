package parser_test

import (
	"fmt"
	"testing"

	"github.com/makramkd/go-monkey/ast"
	"github.com/makramkd/go-monkey/lexer"
	"github.com/makramkd/go-monkey/parser"
	"github.com/makramkd/go-monkey/token"
	"github.com/stretchr/testify/assert"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
let helloWorld = 5 * 3 + 4;
`
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 4)
	assert.Empty(t, p.Errors())

	testCases := []struct {
		expectedIdentifier string
		expectedValue      ast.Expression
	}{
		{"x", &ast.IntegerLiteral{Token: token.New(token.INT, "5"), Value: int64(5)}},
		{"y", &ast.IntegerLiteral{Token: token.New(token.INT, "10"), Value: int64(10)}},
		{"foobar", &ast.IntegerLiteral{Token: token.New(token.INT, "838383"), Value: int64(838383)}},
		{"helloWorld", &ast.InfixExpression{
			Token: token.New(token.PLUS, "+"),
			Left: &ast.InfixExpression{
				Token: token.New(token.TIMES, "*"),
				Left: &ast.IntegerLiteral{
					Token: token.New(token.INT, "5"),
					Value: int64(5),
				},
				Operator: "*",
				Right: &ast.IntegerLiteral{
					Token: token.New(token.INT, "3"),
					Value: int64(3),
				},
			},
			Operator: "+",
			Right: &ast.IntegerLiteral{
				Token: token.New(token.INT, "4"),
				Value: int64(4),
			},
		},
		},
	}

	for i, testCase := range testCases {
		stmt := program.Statements[i]
		testLetStatement(t, stmt, testCase.expectedIdentifier)
		letStmt := stmt.(*ast.LetStatement)
		assert.IsType(t, testCase.expectedValue, letStmt.Value)
		assert.Equal(t, testCase.expectedValue, letStmt.Value)
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

func TestIntegerLiteral(t *testing.T) {
	input := `5;`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	assert.Empty(t, p.Errors())
	assert.Len(t, program.Statements, 1)
	assert.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	assert.IsType(t, &ast.IntegerLiteral{}, stmt.Expression)

	literal := stmt.Expression.(*ast.IntegerLiteral)
	assert.Equal(t, int64(5), literal.Value)
	assert.Equal(t, "5", literal.TokenLiteral())
}

func TestPrefixExpressions(t *testing.T) {
	testCases := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		assert.Empty(t, p.Errors())
		assert.Len(t, program.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		assert.IsType(t, &ast.PrefixExpression{}, stmt.Expression)
		exp := stmt.Expression.(*ast.PrefixExpression)
		assert.Equal(t, testCase.operator, exp.Operator)
		testIntegerLiteral(t, exp.Right, testCase.integerValue)
	}
}

func testIntegerLiteral(t *testing.T, right ast.Expression, value int64) {
	assert.IsType(t, &ast.IntegerLiteral{}, right)
	il := right.(*ast.IntegerLiteral)
	assert.Equal(t, value, il.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), il.TokenLiteral())
}

func TestInfixExpressions(t *testing.T) {
	testCases := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 % 5;", 5, "%", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		assert.Empty(t, p.Errors())
		assert.Len(t, program.Statements, 1)

		assert.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		assert.IsType(t, &ast.InfixExpression{}, stmt.Expression)

		exp := stmt.Expression.(*ast.InfixExpression)

		testIntegerLiteral(t, exp.Left, testCase.leftValue)
		assert.Equal(t, testCase.operator, exp.Operator)
		testIntegerLiteral(t, exp.Right, testCase.rightValue)
	}
}

func TestPrecedence(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"!++--a", "(!(++(--a)))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a / b / c", "((a / b) / c)"},
		{"a * b / c", "((a * b) / c)"},
		{`a % b / c`, `((a % b) / c)`},
		{"a * b + c", "((a * b) + c)"},
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		assert.Empty(t, p.Errors())
		actual := program.String()
		assert.Equal(t, testCase.expected, actual)
	}
}
