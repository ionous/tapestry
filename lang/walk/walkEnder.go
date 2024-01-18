package walk

// signature for OnEnd
type Ender func(Walker, error) error

// after evt.End() run a function.
// the function receives the error that the End function returned.
// the callback can "upgrade" errors ( from no error to some error )
// but otherwise the callback should return the error it received.
func OnEnd(evt Events, cb Ender) Events {
	return ender{evt, cb}
}

type ender struct {
	Events
	cb Ender
}

func (n ender) End(w Walker) error {
	e := n.Events.End(w)
	return n.cb(w, e)
}
