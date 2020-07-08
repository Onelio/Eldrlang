package object

type Context struct {
	store  map[string]Object
	parent *Context
}

func NewContext() *Context {
	s := make(map[string]Object)
	return &Context{store: s}
}

func (e *Context) Get(name string) Object {
	obj, _ := e.store[name]
	return obj
}

func (e *Context) Set(name string, val Object) {
	e.store[name] = val
}
