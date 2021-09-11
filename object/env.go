package object

type Env struct {
	store map[string]Object
	outer *Env
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
