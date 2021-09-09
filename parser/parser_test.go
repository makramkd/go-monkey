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
		input    string
		operator string
		value    ast.Expression
	}{
		{"!5;", "!", &ast.IntegerLiteral{Token: token.New(token.INT, "5"), Value: int64(5)}},
		{"-15;", "-", &ast.IntegerLiteral{Token: token.New(token.INT, "15"), Value: int64(15)}},
		{"!true;", "!", &ast.BooleanLiteral{Token: token.New(token.TRUE, "true"), Value: true}},
		{"!false;", "!", &ast.BooleanLiteral{Token: token.New(token.FALSE, "false"), Value: false}},
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
		assert.Equal(t, exp.Right, testCase.value)
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
		{"a + (b + c)", "(a + (b + c))"},
		{"!(true == true)", "(!(true == true))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"a ** 2 * 2 + 1", "(((a ** 2) * 2) + 1)"},
		{"a ** add(1, 2, 3) + 4", "((a ** add(1,2,3)) + 4)"},
		{"a && b || c", "((a && b) || c)"},
		{"a || b && c", "(a || (b && c))"},
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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x + y; y - x; x**2; } else { x + y; }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	assert.Empty(t, p.Errors())
	assert.Len(t, program.Statements, 1)

	assert.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assert.IsType(t, &ast.IfExpression{}, stmt.Expression)
	ifExp := stmt.Expression.(*ast.IfExpression)

	expectedCondition := &ast.InfixExpression{
		Token: token.New(token.LESS_THAN, "<"),
		Left: &ast.Identifier{
			Token: token.New(token.IDENT, "x"),
			Value: "x",
		},
		Operator: "<",
		Right: &ast.Identifier{
			Token: token.New(token.IDENT, "y"),
			Value: "y",
		},
	}

	expectedConsequence := &ast.BlockStatement{
		Token: token.New(token.LBRACE, "{"),
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Token: token.New(token.IDENT, "x"),
				Expression: &ast.InfixExpression{
					Token: token.New(token.PLUS, "+"),
					Left: &ast.Identifier{
						Token: token.New(token.IDENT, "x"),
						Value: "x",
					},
					Operator: "+",
					Right: &ast.Identifier{
						Token: token.New(token.IDENT, "y"),
						Value: "y",
					},
				},
			},
			&ast.ExpressionStatement{
				Token: token.New(token.IDENT, "y"),
				Expression: &ast.InfixExpression{
					Token: token.New(token.MINUS, "-"),
					Left: &ast.Identifier{
						Token: token.New(token.IDENT, "y"),
						Value: "y",
					},
					Operator: "-",
					Right: &ast.Identifier{
						Token: token.New(token.IDENT, "x"),
						Value: "x",
					},
				},
			},
			&ast.ExpressionStatement{
				Token: token.New(token.IDENT, "x"),
				Expression: &ast.InfixExpression{
					Token: token.New(token.POWER, "**"),
					Left: &ast.Identifier{
						Token: token.New(token.IDENT, "x"),
						Value: "x",
					},
					Operator: "**",
					Right: &ast.IntegerLiteral{
						Token: token.New(token.INT, "2"),
						Value: int64(2),
					},
				},
			},
		},
	}

	expectedAlternative := &ast.BlockStatement{
		Token: token.New(token.LBRACE, "{"),
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Token: token.New(token.IDENT, "x"),
				Expression: &ast.InfixExpression{
					Token: token.New(token.PLUS, "+"),
					Left: &ast.Identifier{
						Token: token.New(token.IDENT, "x"),
						Value: "x",
					},
					Operator: "+",
					Right: &ast.Identifier{
						Token: token.New(token.IDENT, "y"),
						Value: "y",
					},
				},
			},
		},
	}

	assert.Equal(t, expectedCondition, ifExp.Condition)
	assert.Equal(t, expectedConsequence, ifExp.Consequence)
	assert.Equal(t, expectedAlternative, ifExp.Alternative)
}

func TestFunctionLiteral(t *testing.T) {
	input := `fn (x, y) { return x + y; }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	assert.Empty(t, p.Errors())
	assert.Len(t, program.Statements, 1)
	assert.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assert.IsType(t, &ast.FunctionLiteral{}, stmt.Expression)
	fLit := stmt.Expression.(*ast.FunctionLiteral)

	expectedParameters := []*ast.Identifier{
		{
			Token: token.New(token.IDENT, "x"),
			Value: "x",
		},
		{
			Token: token.New(token.IDENT, "y"),
			Value: "y",
		},
	}

	expectedBody := &ast.BlockStatement{
		Token: token.New(token.LBRACE, "{"),
		Statements: []ast.Statement{
			&ast.ReturnStatement{
				Token: token.New(token.RETURN, "return"),
				ReturnValue: &ast.InfixExpression{
					Token: token.New(token.PLUS, "+"),
					Left: &ast.Identifier{
						Token: token.New(token.IDENT, "x"),
						Value: "x",
					},
					Operator: "+",
					Right: &ast.Identifier{
						Token: token.New(token.IDENT, "y"),
						Value: "y",
					},
				},
			},
		},
	}

	assert.Equal(t, expectedParameters, fLit.Parameters)
	assert.Equal(t, expectedBody, fLit.Body)
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2, 3);`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	assert.Empty(t, p.Errors())
	assert.Len(t, program.Statements, 1)

	expectedCall := &ast.CallExpression{
		Token: token.New(token.LPAREN, "("),
		Function: &ast.Identifier{
			Token: token.New(token.IDENT, "add"),
			Value: "add",
		},
		Arguments: []ast.Expression{
			&ast.IntegerLiteral{
				Token: token.New(token.INT, "1"),
				Value: int64(1),
			},
			&ast.IntegerLiteral{
				Token: token.New(token.INT, "2"),
				Value: int64(2),
			},
			&ast.IntegerLiteral{
				Token: token.New(token.INT, "3"),
				Value: int64(3),
			},
		},
	}
	expectedExprStmt := &ast.ExpressionStatement{
		Token:      token.New(token.IDENT, "add"),
		Expression: expectedCall,
	}

	assert.Equal(t, expectedExprStmt, program.Statements[0])
}
