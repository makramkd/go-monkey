package object

import (
	"fmt"
	"strconv"
)

type ObjectType string

const (
	INTEGER ObjectType = "INTEGER"
	BOOLEAN ObjectType = "BOOLEAN"
	NULL    ObjectType = "NULL"
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
