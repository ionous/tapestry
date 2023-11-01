package charm

type State interface {
	// process the next element of the incoming data,
	// return the next state or nil when done.
	NewRune(rune) State
}

// nil represents unhandled runes
// so this function always returns nil
func Unhandled() State { return nil }

// a next state indicating an terminal error
func Error(e error) State {
	return terminalState{err: e}
}

// a next state after which nothing else can be parsed
// but is not otherwise in error.
func Quit() State {
	return terminalState{err: nil}
}

// for the very next rune, return nil ( unhandled )
// it may be the end of parsing, or some parent state might be taking over from here on out.
func Exit(reason string) State {
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
