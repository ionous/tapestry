package charm

// RunStep - run a sequence of of two states
// sending the current rune to the first state immediately
// see also: Step()
func RunStep(r rune, first, last State) State {
	return Step(first, last).NewRune(r)
}

// Step - construct a sequence of two states.
// If the next rune is not handled by the first state or any of its returned states,
// the rune is handed to the second state.
// this acts similar to a parent-child statechart.
func Step(first, last State) State {
	return &chainParser{first, last}
}

// For use in Step() to run an action after the first step completes.
func OnExit(name string, onExit func()) State {
	return Statement("on exit", func(rune) (none State) {
		onExit()
		return
	})
}

type chainParser struct {
	next, last State
}

func (p *chainParser) String() string {
	return StateName(p.next) + "(chain: " + StateName(p.last) + ")"
}

// runs the each state ( and any of their returned states ) to completion
func (p *chainParser) NewRune(r rune) (ret State) {
	if next := p.next.NewRune(r); next == nil {
		// out of next states, run the original last state
		ret = p.last.NewRune(r)
	} else if err, ok := next.(Terminal); ok {
		// if the next state is an error state, return it now.
		ret = err
	} else {
		// remember the new next state, and
		// return *this* to keep stepping towards last.
		ret, p.next = p, next
	}
	return
}
