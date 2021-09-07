package parser

import (
	"fmt"

	"github.com/makramkd/go-monkey/ast"
	"github.com/makramkd/go-monkey/lexer"
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
	EQUALS      // ==
	LESSGREATER // >, >=, <, or <=
	SUM         // + or -
	PRODUCT     // * or /
	PREFIX      // -X or !X
	CALL        // function(X)
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
	errors    []error

	prefixParseFuncs map[token.Type]prefixParseFunc
	infixParseFuncs  map[token.Type]infixParseFunc
}

func New(l *lexer.Lexer) *Parser {
	parser := &Parser{
		l:                l,
		errors:           []error{},
		prefixParseFuncs: map[token.Type]prefixParseFunc{},
		infixParseFuncs:  map[token.Type]infixParseFunc{},
	}

	parser.registerPrefix(token.IDENT, parser.parseIdentifier)

	// Read two tokens so that curToken and peekToken are set
	parser.nextToken()
	parser.nextToken()

	return parser
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for p.curToken.T != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.T {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// Semicolons are optional here
	// TODO: why?
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStmt := &ast.LetStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	letStmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// letStmt.Value = p.parseExpression()
	// TODO: skipping expressions for now
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return letStmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{
		Token: p.curToken,
	}

	// returnStmt.ReturnValue = p.parseExpression()
	// TODO: skipping expressions for now
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return returnStmt
}

func (p *Parser) parseExpression(precedence operatorPrecedence) ast.Expression {
	prefix := p.prefixParseFuncs[p.curToken.T]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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
		fmt.Errorf("expected next token to be %s, got %s instead", t, p.peekToken.T))
}

func (p *Parser) registerPrefix(t token.Type, f prefixParseFunc) {
	p.prefixParseFuncs[t] = f
}

func (p *Parser) registerInfix(t token.Type, f infixParseFunc) {
	p.infixParseFuncs[t] = f
}
