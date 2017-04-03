package object

func NewScope(p *Scope) *Scope {
	s := make(map[string]Object)
	return &Scope{store: s, parentScope: p}
}

type Scope struct {
	store       map[string]Object
	parentScope *Scope
}

func (s *Scope) Get(name string) (Object, bool) {
	obj, ok := s.store[name]
	if !ok && s.parentScope != nil {
		obj, ok = s.parentScope.Get(name)
	}
	return obj, ok
}

func (s *Scope) Set(name string, val Object) Object {
	s.store[name] = val
	return val
}
