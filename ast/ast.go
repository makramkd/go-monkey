package ast

import (
	"strings"

	"github.com/makramkd/go-monkey/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program represents a Monkey program.
// Monkey programs are a sequence of statements that are executed in the order
// in which they are written.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	builder := strings.Builder{}

	for _, s := range p.Statements {
		builder.WriteString(s.String())
	}

	return builder.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string { return i.Value }

// LetStatement represents a Monkey let statement.
// e.g let a = b;
// A let statement consists of an identifier, which appears on the left hand side
// of the assignment operator, and an expression, which appears on the right hand side
// of the assignment operator.
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) statementNode() {}

func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

func (l *LetStatement) String() string {
	builder := strings.Builder{}
	builder.WriteString(l.TokenLiteral() + " ")
	builder.WriteString(l.Name.String())
	builder.WriteString(" = ")

	if l.Value != nil {
		builder.WriteString(l.Value.String())
	}

	builder.WriteRune(';')

	return builder.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	builder := strings.Builder{}
	builder.WriteString(r.TokenLiteral() + " ")

	if r.ReturnValue != nil {
		builder.WriteString(r.ReturnValue.String())
	}

	builder.WriteRune(';')
	return builder.String()
}

type ExpressionStatement struct {
	// The first token of the expression
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}
