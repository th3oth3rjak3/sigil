package interpreter

type EvaluatorEnvironment struct {
	store map[string]Object
	outer *EvaluatorEnvironment
}

func NewEnclosedEvaluatorEnvironment(outer *EvaluatorEnvironment) *EvaluatorEnvironment {
	env := NewEvaluatorEnvironment()
	env.outer = outer
	return env
}

func NewEvaluatorEnvironment() *EvaluatorEnvironment {
	s := make(map[string]Object)
	return &EvaluatorEnvironment{store: s}
}

func (e *EvaluatorEnvironment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *EvaluatorEnvironment) Contains(name string) bool {
	_, ok := e.store[name]
	return ok
}

func (e *EvaluatorEnvironment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}
