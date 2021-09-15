package object

type ExecutionContext int

const (
	ExecutionContextNone ExecutionContext = iota
	ExecutionContextLoop
)

type Env struct {
	store            map[string]Object
	outer            *Env
	executionContext ExecutionContext
}

func NewEnv() *Env {
	return &Env{
		store: map[string]Object{},
		outer: nil,
	}
}

func NewScopedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
}

func (e *Env) Get(name string) (o Object, ok bool) {
	o, ok = e.store[name]
	if !ok && e.outer != nil {
		o, ok = e.outer.Get(name)
	}
	return
}

func (e *Env) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}

func (e *Env) SetExecutionContext(context ExecutionContext) {
	e.executionContext = context
}

func (e *Env) GetExecutionContext() ExecutionContext {
	return e.executionContext
}
