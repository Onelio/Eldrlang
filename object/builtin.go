package object

type BuiltFun func(args ...Object) Object

var Builtins = map[string]*Builtin{
	"len": {
		Size: 1,
		Fun: func(args ...Object) Object {
			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return nil
			}
		},
	},
}
