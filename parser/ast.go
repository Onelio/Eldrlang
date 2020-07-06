package parser

import (
	"Eldrlang/util"
	"bytes"
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

type Package struct {
	Node
	Nodes  []Node
	Errors util.Errors
}

func (p *Package) String() string {
	var out bytes.Buffer
	var i int
	for _, line := range p.Nodes {
		switch line.(type) {
		case *Block, *Conditional, *Function, *Loop:
			lines := strings.Split(line.String(), "\n")
			for _, sub := range lines {
				_, _ = fmt.Fprintf(&out, "%d\t%s\n", i, sub)
				i++
			}
		default:
			_, _ = fmt.Fprintf(&out, "%d\t%s;\n", i, line.String())
			i++
		}
	}
	return out.String()
}
