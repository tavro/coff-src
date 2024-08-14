package object

import (
	"fmt"
)

type ObjectType string

const (
	INT_OBJ = "INT"
	BOOL_OBJ = "BOOL"
	NULL_OBJ = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Int struct {
	Value int64
}

type Bool struct {
	Value bool
}

type Null struct {

}

func (i *Int) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Int) Type() ObjectType { return INT_OBJ }

func (b *Bool) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *Bool) Type() ObjectType { return BOOL_OBJ }

func (n *Null) Inspect() string { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }