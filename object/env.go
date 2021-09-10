package object

type Env struct {
	store map[string]Object
}

func NewEnv() *Env {
	return &Env{
		store: map[string]Object{},
	}
}

func (e *Env) Get(name string) (o Object, ok bool) {
	o, ok = e.store[name]
	return
}

func (e *Env) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}
