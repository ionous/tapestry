package walk

import "errors"

type Events struct {
	BeforeSlot, AfterSlot, OnCommand Event
}

// if next is not nil, then next becomes the next handler.
// ( next is ignored for AfterSlot )
type Event func(w Walker) (next Events, err error)

// can be used to early exit from visiting
// caller still has to check for this to disambiguate the returned error
var DoneVisiting = errors.New("done visiting")

// turn the iterator into an event style traversal
// ( helps especially with backwards compatibility )
func VisitCommand(w Walker, evts Events) error {
	return evts.visitFlow(w)
	// fix: if we looked at the outer type instead of the fields
	// this would be a lot more flexible
	// currently, requires a block
}

func VisitSlot(w Walker, evts Events) error {
	return evts.visitSlot(w)
}

func merge(prev, next Events) Events {
	if next.BeforeSlot == nil {
		next.BeforeSlot = prev.BeforeSlot
	}
	if next.AfterSlot == nil {
		next.AfterSlot = prev.AfterSlot
	}
	if next.OnCommand == nil {
		next.OnCommand = prev.OnCommand
	}
	return next
}

func (b Events) visitFields(w Walker) (err error) {
	for w.Next() && err == nil {
		switch f := w.Field(); f.SpecType() {
		case Flow:
			if !f.Repeats() {
				err = b.visitFlow(w.Walk())
			} else {
				err = b.visitFlows(w.Walk())
			}
		case Slot:
			if !f.Repeats() {
				err = b.visitSlot(w.Walk())
			} else {
				err = b.visitSlots(w.Walk())
			}
		}
	}
	return
}

func (b Events) visitFlow(w Walker) (err error) {
	if n, e := call(w, b.OnCommand); e != nil {
		err = e
	} else {
		next := merge(b, n)
		err = next.visitFields(w)
	}
	return
}

func (b Events) visitSlot(w Walker) (err error) {
	if n, e := call(w, b.BeforeSlot); e != nil {
		err = e
	} else {
		if flow, ok := unpackSlot(w); ok {
			b = merge(b, n)
			err = b.visitFlow(flow)
		}
	}
	if err == nil {
		if _, e := call(w, b.AfterSlot); e != nil {
			err = e
		}
	}
	return
}

func (b Events) visitFlows(w Walker) (err error) {
	for w.Next() {
		if e := b.visitFlow(w.Walk()); e != nil {
			err = e
			break
		}
	}
	return
}

func (b Events) visitSlots(w Walker) (err error) {
	for w.Next() {
		if e := b.visitSlot(w.Walk()); e != nil {
			err = e
			break
		}
	}
	return
}

// returns done
func call(w Walker, cb Event) (ret Events, err error) {
	if cb != nil {
		if n, e := cb(w); e != nil {
			err = e
		} else {
			ret = n
		}
	}
	return
}

// given a slot, return its flow ( false if the slot was empty )
func unpackSlot(slot Walker) (ret Walker, okay bool) {
	if slot.Next() { // next moves the focus to the command if any
		ret, okay = slot.Walk(), true // walk returns the command as a container
	}
	return
}
