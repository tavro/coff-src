package object

import (
	"fmt"
	"coff-src/src/coff/ast"
	"strings"
	"bytes"
)

type ObjectType string

const (
	INT_OBJ = "INT"
	BOOL_OBJ = "BOOL"
	NULL_OBJ = "NULL"
	RET_VAL_OBJ = "RET_VAL"
	ERR_OBJ = "ERR"
	FUN_OBJ = "FUN"
	STR_OBJ = "STR"
	STD_OBJ = "STD"
	ARR_OBJ = "ARR"
)

type Arr struct {
	Elements []Object
}

type StdFunction func(args ...Object) Object
type Std struct {
	Fun StdFunction
}

type Function struct {
	Parameters []*ast.Identifier
	Body *ast.BlockStatement
	Env *Env
}

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

type Str struct {
	Value string
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

func (f *Function) Type() ObjectType { return FUN_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func (s *Str) Type() ObjectType { return STR_OBJ }
func (s *Str) Inspect() string { return s.Value }

func (std *Std) Type() ObjectType { return STD_OBJ }
func (std *Std) Inspect() string { return "std function" }

func (a *Arr) Type() ObjectType { return ARR_OBJ }
func (a *Arr) Inspect() string {
	var out bytes.Buffer
	elements := []string{}

	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}