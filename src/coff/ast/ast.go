package ast

import (
	"coff-src/src/coff/token"
	"bytes"
)

type Boolean struct {
	Token token.Token
	Value bool
}

type InfixExpression struct {
	Token token.Token
	Left Expression
	Operator string
	Right Expression
}

type PrefixExpression struct {
	Token token.Token
	Operator string
	Right Expression
}

type IntLiteral struct {
	Token token.Token
	Value int64
}

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

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (ds DefStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ds.TokenLiteral() + " ")
	out.WriteString(ds.Name.String())
	out.WriteString(" = ")

	if ds.Value != nil {
		out.WriteString(ds.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (rs *RetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	
	if rs.RetValue != nil {
		out.WriteString(rs.RetValue.String())
	}

	out.WriteString(";")
	
	return out.String()
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	
	return ""
}

func (i *Identifier) String() string { return i.Value }

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type DefStatement struct {
	Token token.Token
	Name *Identifier
	Value Expression
}

func (ds *DefStatement) statementNode() {}
func (ds *DefStatement) TokenLiteral() string { return ds.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type RetStatement struct {
	Token token.Token
	RetValue Expression
}

func (rs *RetStatement) statementNode() {}
func (rs *RetStatement) TokenLiteral() string { return rs.Token.Literal }

type ExpressionStatement struct {
	Token token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (il *IntLiteral) expressionNode() {}
func (il *IntLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntLiteral) String() string { return il.Token.Literal }

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

func (oe *InfixExpression) expressionNode() {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string { return b.Token.Literal }

