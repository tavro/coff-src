package object

import (
	"fmt"
)

type ObjectType string

const (
	INT_OBJ = "INT"
	BOOL_OBJ = "BOOL"
	NULL_OBJ = "NULL"
	RET_VAL_OBJ = "RET_VAL"
	ERR_OBJ = "ERR"
)

type Error struct {
	Message string
}

type RetVal struct {
	Value Object
}

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

func (rv *RetVal) Type() ObjectType { return RET_VAL_OBJ }
func (rv *RetVal) Inspect() string { return rv.Value.Inspect() }

func (e *Error) Type() ObjectType { return ERR_OBJ }
func (e *Error) Inspect() string { return "ERROR: " + e.Message }