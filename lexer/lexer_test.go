package lexer_test

import (
	"testing"

	"github.com/makramkd/go-monkey/lexer"
	"github.com/makramkd/go-monkey/token"
	"github.com/stretchr/testify/assert"
)

func TestNextTokenBasic(t *testing.T) {
	input := `=+(){},;`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	lex := lexer.New(input)
	for _, tt := range tests {
		tok := lex.NextToken()

		assert.Equal(t, tt.expectedLiteral, tok.Literal)
		assert.Equal(t, tt.expectedType, tok.T)
	}
}

func TestNextTokenCode(t *testing.T) {
	input := `
let five = 5;
let ten = 10;
let two = 10 / 5;
let mightBeTrue = 4 > 2;
let mightBeFalse = 4 < 2;

let add = fn(x, y) { x + y; };

let result = add(five, ten);

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10
10 != 9

a--
b++
c ** 2
d %= 2

a -= 1
a += 1
a *= 1
a /= 1
`
	l := lexer.New(input)
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		// let five = 5;
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		// let ten = 10;
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		// let two = 10 / 5;
		{token.LET, "let"},
		{token.IDENT, "two"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.DIVIDE, "/"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		// let mightBeTrue = 4 > 2;
		{token.LET, "let"},
		{token.IDENT, "mightBeTrue"},
		{token.ASSIGN, "="},
		{token.INT, "4"},
		{token.GREATER_THAN, ">"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},

		// let mightBeFalse = 4 < 2;
		{token.LET, "let"},
		{token.IDENT, "mightBeFalse"},
		{token.ASSIGN, "="},
		{token.INT, "4"},
		{token.LESS_THAN, "<"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},

		// let add = fn(x, y) { x + y; };
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		// let result = add(five, ten);
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		// if (5 < 10) {
		// 	return true;
		// } else {
		// 	return false;
		// }
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LESS_THAN, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		// 10 == 10
		{token.INT, "10"},
		{token.EQUAL, "=="},
		{token.INT, "10"},

		// 10 != 9
		{token.INT, "10"},
		{token.NOT_EQUAL, "!="},
		{token.INT, "9"},

		// a--
		{token.IDENT, "a"},
		{token.DECR_ONE, "--"},

		// b++
		{token.IDENT, "b"},
		{token.INCR_ONE, "++"},

		// c ** 2
		{token.IDENT, "c"},
		{token.POWER, "**"},
		{token.INT, "2"},

		// d %= 2
		{token.IDENT, "d"},
		{token.REM_EQ, "%="},
		{token.INT, "2"},

		// a -= 1
		{token.IDENT, "a"},
		{token.DECR, "-="},
		{token.INT, "1"},

		// a += 1
		{token.IDENT, "a"},
		{token.INCR, "+="},
		{token.INT, "1"},

		// a *= 1
		{token.IDENT, "a"},
		{token.TIMES_EQ, "*="},
		{token.INT, "1"},

		// a /= 1
		{token.IDENT, "a"},
		{token.DIV_EQ, "/="},
		{token.INT, "1"},
	}

	for _, tt := range tests {
		tok := l.NextToken()

		assert.Equal(t, tt.expectedType, tok.T)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}
