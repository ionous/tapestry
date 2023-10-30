package charm

type actionState struct {
	name    string
	closure func()
}

func (s *actionState) StateName() string {
	return s.name
}

// NewRune implements State by calling the action and returning nil.
func (s *actionState) NewRune(r rune) State {
	s.closure()
	return nil
}

// StateExit provides a callback when a state ends
func StateExit(name string, onExit func()) State {
	return &actionState{"exit " + name, onExit}
}

// for the very next rune, return nil ( unhandled )
var Terminal = Statement("terminal", func(rune) State {
	return nil
})
