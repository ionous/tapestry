package charm

// represents both the error state and the final state
type terminalState struct {
	err error
}

func (e terminalState) String() (ret string) {
	if e.err != nil {
		ret = "error state"
	} else {
		ret = "terminal state"
	}
	return
}

func (e terminalState) NewRune(r rune) (ret State) {
	return e
}
