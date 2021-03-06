package parser

import (
	"fmt"
	"strconv"

	"github.com/makramkd/go-monkey/ast"
	"github.com/makramkd/go-monkey/lexer"
	"github.com/makramkd/go-monkey/token"
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

	parser.registerPrefixes()
	parser.registerInfixes()

	// Read two tokens so that curToken and peekToken are set
	parser.nextToken()
	parser.nextToken()

	return parser
}

func (p *Parser) Errors() []error {
	return p.errors
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
	case token.IMPORT:
		return p.parseImportStatement()
	case token.FOR:
		return p.parseForEachStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: p.curToken}

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseForEachStatement() *ast.ForEachStatement {
	stmt := &ast.ForEachStatement{Token: p.curToken}

	stmt.Identifiers = p.parseForIdentifiers()

	// The 'in' token is required.
	if !p.expectPeek(token.IN) {
		return nil
	}

	p.nextToken()

	stmt.Collection = p.parseForCollection()

	p.nextToken()

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForIdentifiers() []*ast.Identifier {
	// Comma separated list of identifiers, similar to a function.
	ids := []*ast.Identifier{}

	// Need at least one identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	ids = append(ids, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if !p.expectPeek(token.IDENT) {
			return nil
		}

		ids = append(ids, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	return ids
}

func (p *Parser) parseForCollection() ast.Expression {
	switch p.curToken.T {
	case token.LBRACE:
		return p.parseHashLiteral()
	case token.LBRACK:
		return p.parseArrayLiteral()
	case token.IDENT:
		return p.parseIdentifier()
	default:
		return nil
	}
}

func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Module = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return stmt
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

	// Advance to the next token because the assignment operator doesn't
	// have prefix or infix expression methods.
	p.nextToken()

	letStmt.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return letStmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{
		Token: p.curToken,
	}

	p.nextToken()
	returnStmt.ReturnValue = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return returnStmt
}

func (p *Parser) parseExpression(precedence operatorPrecedence) ast.Expression {
	prefix := p.prefixParseFuncs[p.curToken.T]
	if prefix == nil {
		p.noPrefixParseFuncError(p.curToken.T)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFuncs[p.peekToken.T]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("could not parse %q as integer: %v", p.curToken.Literal, err))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	lit := &ast.BooleanLiteral{Token: p.curToken}

	value, err := strconv.ParseBool(p.curToken.Literal)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("could not parse %q as bool: %v", p.curToken.Literal, err))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	lit := &ast.ArrayLiteral{Token: p.curToken}

	lit.Elements = p.parseArrayElements()

	return lit
}

func (p *Parser) parseArrayElements() []ast.Expression {
	elements := []ast.Expression{}

	// Check for empty array
	if p.peekTokenIs(token.RBRACK) {
		p.nextToken()
		return elements
	}

	p.nextToken()
	elements = append(elements, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		elements = append(elements, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RBRACK) {
		return nil
	}

	return elements
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}

	p.nextToken()

	// Recursively parse the expression past the operator token
	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken, Statements: []ast.Statement{}}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	parameters := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return parameters
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	parameters = append(parameters, &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		parameters = append(parameters, &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return parameters
}

func (p *Parser) parseIndexAccessExpression(array ast.Expression) ast.Expression {
	accessExpr := &ast.IndexAccessExpression{
		Token: p.curToken,
		Left:  array,
	}

	p.nextToken()

	accessExpr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACK) {
		return nil
	}

	return accessExpr
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	callExpr := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	callExpr.Arguments = p.parseCallArguments()
	return callExpr
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseHashLiteral() ast.Expression {
	lit := &ast.HashLiteral{Token: p.curToken}

	pairs := p.parseHashPairs()
	if pairs == nil {
		return nil
	}

	lit.Pairs = pairs

	return lit
}

func (p *Parser) parseHashPairs() map[ast.Expression]ast.Expression {
	pairs := map[ast.Expression]ast.Expression{}

	// handle empty hash
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return pairs
	}

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return pairs
}
