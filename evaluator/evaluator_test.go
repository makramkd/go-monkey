package evaluator_test

import (
	"testing"

	"github.com/makramkd/go-monkey/evaluator"
	"github.com/makramkd/go-monkey/lexer"
	"github.com/makramkd/go-monkey/object"
	"github.com/makramkd/go-monkey/parser"
	"github.com/stretchr/testify/assert"
)

func TestEvalIntegerLiteral(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 4", 9},
		{"5 + 4 * 10", 45},
		{"28 / 2 + 3 * 4 + 1", 27},
		{"(4 + 10) * 2 + (3 + 10) * 2 + 1", 55},
		{"2 ** 2", 4},
		{"2 % 2", 0},
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnv()
		val := evaluator.Eval(program, env)
		assert.IsType(t, &object.Integer{}, val)
		integerValue := val.(*object.Integer)
		assert.Equal(t, testCase.expected, integerValue.Value)
	}
}

func TestEvalBooleanLiteral(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 != 2", true},
		{"1 == 1", true},
		{"1 == 2", false},
		{"1 > 1", false},
		{"1 >= 1", true},
		{"2 > 1", true},
		{"2 < 1", false},
		{"2 <= 1", false},
		{"1 <= 2", true},
		{"1 > 2 && 2 > 1", false},
		{"true == true", true},
		{"false != true", true},
		{"1 && 0", false},
		{"1 || 0", true},
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnv()
		val := evaluator.Eval(program, env)
		assert.IsType(t, &object.Boolean{}, val)
		boolValue := val.(*object.Boolean)
		assert.Equal(t, testCase.expected, boolValue.Value)
	}
}

func TestBangOperator(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnv()
		val := evaluator.Eval(program, env)
		assert.IsType(t, &object.Boolean{}, val)
		boolValue := val.(*object.Boolean)
		assert.Equal(t, testCase.expected, boolValue.Value)
	}
}

func TestIfExpressions(t *testing.T) {
	testCases := []struct {
		input    string
		expected object.Object
	}{
		{"if (true) { 10 }", &object.Integer{Value: 10}},
		{"if (true || false) { false }", &object.Boolean{Value: false}},
		{"if (1 < 2 && (3 - 4) == -1) { 42 } else { 41 }", &object.Integer{Value: 42}},
		{"if (false) { 41 } else { 42 }", &object.Integer{Value: 42}},
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnv()
		val := evaluator.Eval(program, env)
		assert.Equal(t, testCase.expected, val)
	}
}

func TestReturnStatements(t *testing.T) {
	testCases := []struct {
		input    string
		expected object.Object
	}{
		{"return 10;", &object.Integer{Value: 10}},
		{"return true;", &object.Boolean{Value: true}},
		{"1 + 1; return if (1 == 1) { 42 } else { 43 };", &object.Integer{Value: 42}},
		{"if (10 > 1) { if (10 > 2) { return 10; } return 1; }", &object.Integer{Value: 10}},
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnv()
		val := evaluator.Eval(program, env)
		assert.Equal(t, testCase.expected, val)
	}
}

func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		input    string
		expected *object.Error
	}{
		{"5 + true;", &object.Error{Message: "type mismatch: INTEGER + BOOLEAN"}},
		{"-true;", &object.Error{Message: "unknown operator: -BOOLEAN"}},
		{"true + false", &object.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"}},
		{"if (10 > 1) { if ( 10 > 2 ) { return false + true; } return 42; }", &object.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"}},
		{"if (true + false == 1) { return 42; }", &object.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"}},
		{"if (true == false * 1) { return 42; }", &object.Error{Message: "type mismatch: BOOLEAN * INTEGER"}},
		{"foobar", &object.Error{Message: "identifier not found: foobar"}},
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnv()
		val := evaluator.Eval(program, env)
		assert.Equal(t, testCase.expected, val)
	}
}

func TestLetStatements(t *testing.T) {
	testCases := []struct {
		input    string
		expected object.Object
	}{
		{"let a = 5; a;", &object.Integer{Value: 5}},
		{"let a = 5 + 5; a;", &object.Integer{Value: 10}},
		{"let a = if (5 > 4) { 42 } else { 41 }; a;", &object.Integer{Value: 42}},
		{"let a = 5; let b = a; let c = a + b + 5; c;", &object.Integer{Value: 15}},
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnv()
		val := evaluator.Eval(program, env)
		assert.Equal(t, testCase.expected, val)
	}
}

func TestFunctionObject(t *testing.T) {
	testCases := []struct {
		input    string
		expected object.Object
	}{
		{`let f = fn(x) { return x + 2; }; f(2);`, &object.Integer{Value: 4}},
		{`let f = fn(x, y) { return x**2 + y**2; }; f(2, 2);`, &object.Integer{Value: 8}},
		{`let x = 2; let f = fn(x) { return x ** 2; }; f(3);`, &object.Integer{Value: 9}},
		{
			`let x = 2; 
			 let f = fn(x) { 
				let inner = fn(y) {
					return y ** 2;
			 	};
				return inner(x + 1);
			 };
			 f(3);
			 `, &object.Integer{Value: 16}},
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		assert.Len(t, p.Errors(), 0)
		env := object.NewEnv()
		val := evaluator.Eval(program, env)
		assert.Equal(t, testCase.expected, val)
	}
}

func TestStringOperations(t *testing.T) {
	testCases := []struct {
		input    string
		expected object.Object
	}{
		{`"hello world" + " today";`, &object.String{Value: "hello world today"}},
		{`let firstName = "Makram"; let lastName = "Kamaleddine"; let f = fn (first, last) { return first + " " + last; }; f(firstName, lastName);`,
			&object.String{Value: "Makram Kamaleddine"}},
		{`"hello world" == "hello world";`, &object.Boolean{Value: true}},
		{`"hello world" != "today";`, &object.Boolean{Value: true}},
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		assert.Empty(t, p.Errors())
		env := object.NewEnv()
		val := evaluator.Eval(program, env)
		assert.Equal(t, testCase.expected, val)
	}
}
