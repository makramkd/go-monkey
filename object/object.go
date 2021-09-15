package object

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/makramkd/go-monkey/ast"
)

type ObjectType string

const (
	INTEGER      ObjectType = "INTEGER"
	BOOLEAN      ObjectType = "BOOLEAN"
	NULL         ObjectType = "NULL"
	RETURN_VALUE ObjectType = "RETURN_VALUE"
	ERROR        ObjectType = "ERROR"
	FUNCTION     ObjectType = "FUNCTION"
	STRING       ObjectType = "STRING"
	ARRAY        ObjectType = "ARRAY"
	BUILTIN      ObjectType = "BUILTIN"
	HASH         ObjectType = "HASH"
	BREAK        ObjectType = "BREAK"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return strconv.FormatBool(b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN }

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL }

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }
func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE }

type Error struct {
	Message string
}

func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Env
}

func (f *Function) Type() ObjectType { return FUNCTION }
func (f *Function) Inspect() string {
	builder := strings.Builder{}

	builder.WriteString("fn(")
	for _, param := range f.Parameters {
		builder.WriteString(param.String())
	}
	builder.WriteString(") {\n")
	builder.WriteString(f.Body.String())
	builder.WriteString("\n}")
	return builder.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING }
func (s *String) Inspect() string  { return s.Value }

type Array struct {
	Values []Object
}

func (a *Array) Type() ObjectType { return ARRAY }
func (a *Array) Inspect() string {
	elems := []string{}
	for _, e := range a.Values {
		elems = append(elems, e.Inspect())
	}
	b := strings.Builder{}
	b.WriteByte('[')
	b.WriteString(strings.Join(elems, ","))
	b.WriteByte(']')
	return b.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	F BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: BOOLEAN, Value: value}
}

func (b *Integer) HashKey() HashKey {
	return HashKey{Type: INTEGER, Value: uint64(b.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: STRING, Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH }
func (h *Hash) Inspect() string {
	builder := strings.Builder{}

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+":"+pair.Value.Inspect())
	}

	builder.WriteByte('{')
	builder.WriteString(strings.Join(pairs, ", "))
	builder.WriteByte('}')
	return builder.String()
}

type Break struct{}

func (b *Break) Inspect() string  { return "break" }
func (b *Break) Type() ObjectType { return BREAK }
