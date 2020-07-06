package object

import "fmt"

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
