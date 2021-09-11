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

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g ! or -
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	builder := strings.Builder{}
	builder.WriteRune('(')
	builder.WriteString(p.Operator)
	builder.WriteString(p.Right.String())
	builder.WriteRune(')')
	return builder.String()
}

type InfixExpression struct {
	Token    token.Token // the operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}

func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }

func (i *InfixExpression) String() string {
	builder := strings.Builder{}
	builder.WriteRune('(')
	builder.WriteString(i.Left.String() + " ")
	builder.WriteString(i.Operator + " ")
	builder.WriteString(i.Right.String())
	builder.WriteRune(')')
	return builder.String()
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode() {}

func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }

func (b *BooleanLiteral) String() string { return b.Token.Literal }

type BlockStatement struct {
	Token      token.Token // the '{' token
	Statements []Statement
}

func (b *BlockStatement) statementNode()       {}
func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BlockStatement) String() string {
	builder := strings.Builder{}
	for _, stmt := range b.Statements {
		builder.WriteString(stmt.String())
	}
	return builder.String()
}

type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode() {}

func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }

func (i *IfExpression) String() string {
	builder := strings.Builder{}
	builder.WriteString("if")
	builder.WriteString(i.Condition.String())
	builder.WriteString(" ")
	builder.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		builder.WriteString("else")
		builder.WriteString(i.Alternative.String())
	}

	return builder.String()
}

type FunctionLiteral struct {
	Token      token.Token // the 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode()      {}
func (f *FunctionLiteral) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionLiteral) String() string {
	builder := strings.Builder{}
	builder.WriteString("fn")
	builder.WriteRune('(')
	for i, param := range f.Parameters {
		builder.WriteString(param.String())
		if i < len(f.Parameters)-1 {
			builder.WriteRune(',')
		}
	}
	builder.WriteRune(')')
	builder.WriteRune('{')
	builder.WriteString(f.Body.String())
	builder.WriteRune('}')
	return builder.String()
}

type CallExpression struct {
	Token     token.Token // the '(' token
	Function  Expression  // identifier or function literal
	Arguments []Expression
}

func (c *CallExpression) expressionNode()      {}
func (c *CallExpression) TokenLiteral() string { return c.Token.Literal }
func (c *CallExpression) String() string {
	builder := strings.Builder{}
	builder.WriteByte('(')
	builder.WriteString(c.Function.String())
	builder.WriteRune('(')
	for i, arg := range c.Arguments {
		builder.WriteString(arg.String())
		if i < len(c.Arguments)-1 {
			builder.WriteRune(',')
		}
	}
	builder.WriteRune(')')
	builder.WriteByte(')')
	return builder.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string       { return s.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (a *ArrayLiteral) expressionNode()      {}
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteral) String() string {
	builder := strings.Builder{}
	builder.WriteByte('[')
	for i, e := range a.Elements {
		builder.WriteString(e.String())
		if i < len(a.Elements)-1 {
			builder.WriteByte(',')
		}
	}
	builder.WriteByte(']')
	return builder.String()
}

type ArrayAccessExpression struct {
	Token token.Token // The '[' token
	Array Expression  // can either be an array literal or an identifier
	Index Expression  // The index to access from the array
}

func (a *ArrayAccessExpression) expressionNode()      {}
func (a *ArrayAccessExpression) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayAccessExpression) String() string {
	builder := strings.Builder{}
	builder.WriteByte('(')
	builder.WriteString(a.Array.String())
	builder.WriteByte('[')
	builder.WriteString(a.Index.String())
	builder.WriteByte(']')
	builder.WriteByte(')')
	return builder.String()
}
