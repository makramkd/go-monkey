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
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		val := evaluator.Eval(program)
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
	}

	for _, testCase := range testCases {
		l := lexer.New(testCase.input)
		p := parser.New(l)
		program := p.ParseProgram()
		val := evaluator.Eval(program)
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
		val := evaluator.Eval(program)
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
		val := evaluator.Eval(program)
		assert.Equal(t, testCase.expected, val)
	}
}
