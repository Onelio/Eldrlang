package object

type Runtime struct {
	context *Context
}

func NewRuntime() *Runtime {
	return &Runtime{context: NewContext()}
}

func (r *Runtime) PushChild() {
	child := NewContext()
	child.parent, r.context = r.context, child
}

func (r *Runtime) PopChild() {
	r.context = r.context.parent
}

func (r *Runtime) GetValue(name string) Object {
	var (
		context = r.context
		object  Object
	)
	for context != nil {
		object = context.Get(name)
		if object != nil {
			break
		}
		context = context.parent
	}
	return object
}

func (r *Runtime) SetValue(name string, val Object) {
	r.context.Set(name, val)
}
