package object

import (
	"fmt"
	"github.com/Onelio/Eldrlang/parser"
)

type Object interface {
	Inspect() string
}

type Null struct{}

func (n *Null) Inspect() string { return "null" }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

type String struct {
	Value string
}

func (s *String) Inspect() string { return s.Value }

type Function struct {
	Parameters []*parser.Identifier
	Body       *parser.Block
}

func (f *Function) Inspect() string {
	return "function"
}

type Builtin struct {
	Size int
	Fun  BuiltFun
}

func (b *Builtin) Inspect() string { return "builtin function" }
