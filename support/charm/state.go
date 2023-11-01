package charm

type State interface {
	// process the next element of the incoming data,
	// return the next state or nil when done.
	NewRune(rune) State
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
