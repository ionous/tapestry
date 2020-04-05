package rt

type EmptyStream bool

func (EmptyStream) Remaining() int {
	return 0
}

func (EmptyStream) HasNext() bool {
	return false
}

func (EmptyStream) GetNext(pv interface{}) error {
	panic("Attempted to advance an empty stream.")
}
