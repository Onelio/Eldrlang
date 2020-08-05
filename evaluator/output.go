package evaluator

import (
	"github.com/Onelio/Eldrlang/object"
	"github.com/Onelio/Eldrlang/util"
)

type Output struct {
	object.Object
	Errors util.Errors
}

func (o *Output) String() string {
	if o.Object != nil {
		return o.Object.Inspect() + "\n"
	}
	return ""
}
