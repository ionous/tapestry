package charm

// ParseChain - run a sequence of of two states
// sending the current rune to the first state immediately
// ( and returning that result )
// see also: MakeChain()
func ParseChain(r rune, first, last State) State {
	return MakeChain(first, last).NewRune(r)
}

// MakeChain - construct a sequence of two states.
// If the next rune is not handled by the first state or any of its returned states,
// the rune is handed to the second state.
func MakeChain(first, last State) State {
	return &chainParser{first, last}
}

type chainParser struct {
	next, last State
}

func (p *chainParser) StateName() string {
	return "chain parser ('" + p.next.StateName() + "' '" + p.last.StateName() + "')"
}

// runs the each state ( and any of their returned states ) to completion
func (p *chainParser) NewRune(r rune) (ret State) {
	if next := p.next.NewRune(r); next != nil {
		ret, p.next = p, next
	} else {
		ret = p.last.NewRune(r)
	}
	return
}
