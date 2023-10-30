package charm

type State interface {
	// process the next element of the incoming data,
	// return the next state or nil when done.
	NewRune(rune) State
	// an identifier used for error reporting and debugging
	StateName() string
}
