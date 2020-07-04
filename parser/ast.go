package parser

import (
	"fmt"
	"strings"
)

type Node interface {
	Literal() string
	String() string
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type Program struct {
	Node
	Nodes []Node
}

func (p *Program) String() {
	var i int
	for _, line := range p.Nodes {
		switch line.(type) {
		case *Block, *Conditional, *Function, *Loop:
			lines := strings.Split(line.String(), "\n")
			for _, sub := range lines {
				fmt.Printf("%d\t%s\n", i, sub)
				i++
			}
		default:
			fmt.Printf("%d\t%s;\n", i, line.String())
			i++
		}
	}
}
