package charm

// definition of states n the state chart
//if State implements Stringer StateName() will use it.
type State interface {
	// process the next element of the incoming data,
	// return the next state or nil when done.
	NewRune(rune) State
}

// calls NewRune on state. sometimes this is a more convenient notation.
func RunState(r rune, state State) (ret State) {
	if state != nil {
		ret = state.NewRune(r)
	}
	return
}

// nil represents unhandled runes
// so this function always returns nil
func Unhandled() State { return nil }

// for the very next rune, return nil ( unhandled )
// it may be the end of parsing, or some parent state might be taking over from here on out.
func Finished(reason string) State {
	return Statement(reason, func(rune) (none State) {
		return
	})
}

// replaceable function for printing the name of a state
// by default uses Stringer's String(), if not implemented it returns "unknown state"
// test packages can overwrite with something that uses package reflect if desired.
var StateName = func(n State) (ret string) {
	if s, ok := n.(interface{ String() string }); !ok {
		ret = "unknown state"
	} else {
		ret = s.String()
	}
	return
}
