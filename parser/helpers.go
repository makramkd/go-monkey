package parser

import (
	"fmt"

	"github.com/makramkd/go-monkey/ast"
	"github.com/makramkd/go-monkey/token"
)

// Pratt parsing main idea: association of parsing functions with token types.
// Whenever this token type is encountered, the parsing functions are called to parse the
// appropriate expression and return an AST node that represents it.
// Each token type can have up to two parsing functions associated with it, depending
// on whether the token is found in a prefix or infix position.
// TODO: what about postfix?
type prefixParseFunc func() ast.Expression
type infixParseFunc func(ast.Expression) ast.Expression

type operatorPrecedence int

const (
	_ operatorPrecedence = iota
	LOWEST
	OR          // ||
	AND         // &&
	EQUALS      // ==, !=
	LESSGREATER // >, >=, <, or <=
	SUM         // + or -
	PRODUCT     // *, /, or %
	POWER       // **
	PREFIX      // -X, !X, --X, ++X
	CALL        // function(X)
)

// Precedence of the binary operators.
// In the case where we have a prefix expression or function call expression, it's parsing method should
// ensure that the higher priority is applied.
var precedenceTable = map[token.Type]operatorPrecedence{
	token.EQUAL:     EQUALS,
	token.NOT_EQUAL: EQUALS,

	token.OR:  OR,
	token.AND: AND,

	token.LESS_THAN:    LESSGREATER,
	token.GREATER_THAN: LESSGREATER,
	token.LEQ:          LESSGREATER,
	token.GEQ:          LESSGREATER,

	token.PLUS:  SUM,
	token.MINUS: SUM,

	token.TIMES:     PRODUCT,
	token.DIVIDE:    PRODUCT,
	token.REMAINDER: PRODUCT,

	token.POWER: POWER,

	token.LPAREN: CALL,
}

func (p *Parser) registerPrefixes() {
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	for _, tokType := range []token.Type{token.TRUE, token.FALSE} {
		p.registerPrefix(tokType, p.parseBooleanLiteral)
	}
	for _, tokType := range []token.Type{token.BANG, token.MINUS, token.INCR_ONE, token.DECR_ONE} {
		p.registerPrefix(tokType, p.parsePrefixExpression)
	}
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
}

func (p *Parser) registerInfixes() {
	infixOperators := []token.Type{
		token.EQUAL, token.NOT_EQUAL,
		token.OR, token.AND,
		token.LESS_THAN, token.GREATER_THAN, token.LEQ, token.GEQ,
		token.PLUS, token.MINUS,
		token.TIMES, token.DIVIDE, token.REMAINDER,
		token.POWER,
	}
	for _, tokType := range infixOperators {
		p.registerInfix(tokType, p.parseInfixExpression)
	}
	p.registerInfix(token.LPAREN, p.parseCallExpression)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.T == t
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.T == t
}

func (p *Parser) peekError(t token.Type) {
	p.errors = append(
		p.errors,
		fmt.Errorf("expected next token to be '%s', got '%s' instead", t, p.peekToken.T))
}

func (p *Parser) registerPrefix(t token.Type, f prefixParseFunc) {
	p.prefixParseFuncs[t] = f
}

func (p *Parser) registerInfix(t token.Type, f infixParseFunc) {
	p.infixParseFuncs[t] = f
}

func (p *Parser) noPrefixParseFuncError(t token.Type) {
	p.errors = append(p.errors, fmt.Errorf("no prefix parse function found for '%s'", t))
}

func (p *Parser) peekPrecedence() operatorPrecedence {
	if p, ok := precedenceTable[p.peekToken.T]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() operatorPrecedence {
	if p, ok := precedenceTable[p.curToken.T]; ok {
		return p
	}
	return LOWEST
}
