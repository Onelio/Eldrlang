package object

func NewContext() *Context {
	s := make(map[string]Object)
	return &Context{store: s}
}

type Context struct {
	store map[string]Object
	child *Context
}

func (e *Context) NewChildren() *Context {
	e.child = NewContext()
	return e.child
}

func (e *Context) Get(name string) Object {
	var (
		ctx    = e
		object Object
		valid  bool
	)
	for ctx != nil {
		object, valid = ctx.store[name]
		if valid {
			break
		}
		ctx = ctx.child
	}
	return object
}

func (e *Context) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
