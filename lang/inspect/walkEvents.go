package inspect

import (
	"errors"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// fix? this should probably be passing type info and not walker itself
type Events interface {
	Flow(Iter) error
	Slot(Iter) error
	Repeat(Iter) error
	// called for every member of a flow
	Field(Iter) error
	// called for each value in a flow that exists
	Value(Iter) error // FIX
	// called after each flow, slot, or repeat
	End(Iter) error
}

// can be used to early exit from visiting
// caller still has to check for this to disambiguate the returned error
var DoneVisiting = errors.New("done visiting")

// turn the iterator into an event style callbacks
func Visit(i typeinfo.Inspector, evt Events) error {
	return visitFlow(Walk(i), evt)
}

func visitFields(w Iter, evt Events) (err error) {
	for w.Next() && err == nil {
		if e := evt.Field(w); e != nil {
			err = e
			break
		}
		switch f := w.Term(); f.Type.(type) {
		case *typeinfo.Flow:
			if !f.Repeats {
				err = visitFlow(w.Walk(), evt)
			} else if e := visitFlows(w.Walk(), evt); e != nil {
				err = e
			} else {
				evt.End(w)
			}
		case *typeinfo.Slot:
			if !f.Repeats {
				err = visitSlot(w.Walk(), evt)
			} else if e := visitSlots(w.Walk(), evt); e != nil {
				err = e
			} else {
				evt.End(w)
			}
		case *typeinfo.Num, *typeinfo.Str:
			if !f.Repeats {
				err = evt.Value(w)
			} else if e := visitValues(w.Walk(), evt); e != nil {
				err = e
			} else {
				evt.End(w)
			}
		}
	}
	return
}

func visitFlow(w Iter, evt Events) (err error) {
	if e := evt.Flow(w); e != nil {
		err = e
	} else if e := visitFields(w, evt); e != nil {
		err = e
	} else {
		err = evt.End(w)
	}
	return
}

func visitSlot(w Iter, evt Events) (err error) {
	if e := evt.Slot(w); e != nil {
		err = e
	} else {
		if flow, ok := unpackSlot(w); ok {
			err = visitFlow(flow, evt)
		}
		if err == nil {
			err = evt.End(w)
		}
	}
	return
}

func visitFlows(w Iter, evt Events) (err error) {
	for w.Next() {
		if e := visitFlow(w.Walk(), evt); e != nil {
			err = e
			break
		}
	}
	return
}

func visitSlots(w Iter, evt Events) (err error) {
	for w.Next() {
		if e := visitSlot(w.Walk(), evt); e != nil {
			err = e
			break
		}
	}
	return
}

func visitValues(w Iter, evt Events) (err error) {
	for w.Next() {
		if e := evt.Value(w); e != nil {
			err = e
			break
		}
	}
	return
}

// given a slot, return its flow ( false if the slot was empty )
func unpackSlot(slot Iter) (ret Iter, okay bool) {
	if slot.Next() { // next moves the focus to the flow if any
		ret, okay = slot.Walk(), true // walk returns the flow as a container
	}
	return
}
