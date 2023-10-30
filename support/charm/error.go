package charm

func Error(e error) State {
	return errorState{err: e}
}

type errorState struct {
	err error
}

func (e errorState) StateName() string {
	return "error state"
}

func (e errorState) NewRune(r rune) (ret State) {
	return e
}
