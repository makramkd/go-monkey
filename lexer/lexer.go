package lexer

import (
	"unicode"

	"github.com/makramkd/go-monkey/token"
)

type Lexer struct {
	// The input being processed by the lexer.
	// TODO: should probably be []rune instead?
	input string
	// Current position in input (points to current char)
	position int
	// Current reading position in input (after current char)
	readPosition int
	// Current char under examination
	// TODO: should probably be rune
	ch byte
}

// New creates a new Monkey lexer for the given input.
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case ';':
		tok = token.New(token.SEMICOLON, string(l.ch))
	case ',':
		tok = token.New(token.COMMA, string(l.ch))

	// Operators
	case '=':
		if nextChar := l.peekChar(); nextChar == '=' {
			l.readChar()
			tok = token.New(token.EQUAL, "==")
		} else {
			tok = token.New(token.ASSIGN, string(l.ch))
		}
	case '+':
		if nextChar := l.peekChar(); nextChar == '=' {
			l.readChar()
			tok = token.New(token.INCR, "+=")
		} else if nextChar == '+' {
			l.readChar()
			tok = token.New(token.INCR_ONE, "++")
		} else {
			tok = token.New(token.PLUS, string(l.ch))
		}
	case '-':
		if nextChar := l.peekChar(); nextChar == '=' {
			l.readChar()
			tok = token.New(token.DECR, "-=")
		} else if nextChar == '-' {
			l.readChar()
			tok = token.New(token.DECR_ONE, "--")
		} else {
			tok = token.New(token.MINUS, string(l.ch))
		}
	case '*':
		if nextChar := l.peekChar(); nextChar == '=' {
			l.readChar()
			tok = token.New(token.TIMES_EQ, "*=")
		} else if nextChar == '*' {
			l.readChar()
			tok = token.New(token.POWER, "**")
		} else {
			tok = token.New(token.TIMES, string(l.ch))
		}
	case '%':
		if nextChar := l.peekChar(); nextChar == '=' {
			l.readChar()
			tok = token.New(token.REM_EQ, "%=")
		} else {
			tok = token.New(token.REMAINDER, string(l.ch))
		}
	case '/':
		if nextChar := l.peekChar(); nextChar == '=' {
			l.readChar()
			tok = token.New(token.DIV_EQ, "/=")
		} else {
			tok = token.New(token.DIVIDE, string(l.ch))
		}
	case '<':
		if nextChar := l.peekChar(); nextChar == '=' {
			l.readChar()
			tok = token.New(token.LEQ, "<=")
		} else {
			tok = token.New(token.LESS_THAN, string(l.ch))
		}
	case '>':
		if nextChar := l.peekChar(); nextChar == '=' {
			l.readChar()
			tok = token.New(token.GEQ, ">=")
		} else {
			tok = token.New(token.GREATER_THAN, string(l.ch))
		}
	case '!':
		if nextChar := l.peekChar(); nextChar == '=' {
			l.readChar()
			tok = token.New(token.NOT_EQUAL, "!=")
		} else {
			tok = token.New(token.BANG, string(l.ch))
		}

	// Parens, braces, brackets
	case '{':
		tok = token.New(token.LBRACE, string(l.ch))
	case '}':
		tok = token.New(token.RBRACE, string(l.ch))
	case '(':
		tok = token.New(token.LPAREN, string(l.ch))
	case ')':
		tok = token.New(token.RPAREN, string(l.ch))
	case '[':
		tok = token.New(token.LBRACK, string(l.ch))

	case '.':
		tok = token.New(token.PERIOD, string(l.ch))

	case '&':
		if nextChar := l.peekChar(); nextChar == '&' {
			l.readChar()
			tok = token.New(token.AND, "&&")
		} else {
			tok = token.New(token.ILLEGAL, string(l.ch))
		}

	case '|':
		if nextChar := l.peekChar(); nextChar == '|' {
			l.readChar()
			tok = token.New(token.OR, "||")
		} else {
			tok = token.New(token.ILLEGAL, string(l.ch))
		}
	case '"':
		tok.Literal = l.readStringLiteral()
		tok.T = token.STRING

	case 0:
		tok = token.New(token.EOF, "")
	default:
		if isLetter(rune(l.ch)) {
			tok.Literal = l.readIdentifier()
			tok.T = token.LookupIdent(tok.Literal)
			return tok
		} else if unicode.IsDigit(rune(l.ch)) {
			tok.T = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = token.New(token.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	return l.read(isLetter)
}

func (l *Lexer) readNumber() string {
	return l.read(unicode.IsDigit)
}

func (l *Lexer) read(cond func(rune) bool) string {
	position := l.position
	for cond(rune(l.ch)) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readStringLiteral() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}
