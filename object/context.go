package object

func NewEnclosedContext(outer *Context) *Context {
	env := NewContext()
	env.outer = outer
	return env
}

func NewContext() *Context {
	s := make(map[string]Object)
	return &Context{store: s, outer: nil}
}

type Context struct {
	store map[string]Object
	outer *Context
}

func (e *Context) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Context) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
