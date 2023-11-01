package charm

// Statement functions behave as a State.
// ( Self works the same, except that passes the state to its function as "self" and this does not )
func Statement(name string, closure func(rune) State) State {
	return &funcState{name, closure}
}

type funcState struct {
	name    string
	closure func(rune) State
}

func (s *funcState) String() string {
	return s.name
}

// NewRune implements State by calling the Statement's underlying function.
func (s *funcState) NewRune(r rune) State {
	return s.closure(r)
}
