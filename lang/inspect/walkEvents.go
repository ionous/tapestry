package inspect

import (
	"errors"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// fix? this should probably be passing type info and not walker itself
type Events interface {
	Flow(Iter) error
	// called for every member of a flow.
	Field(Iter) error
	// called for a member slot: can be an empty slot.
	Slot(Iter) error
	// called for a member that repeats; can be an empty list.
	Repeat(Iter) error
	// called for each str or num in a flow, or in a repeat.
	// because it isnt a container, there is no matching end.
	Value(Iter) error
	// called after each flow, slot, or repeat.
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
			} else {
				err = visitFlows(w.Walk(), evt)
			}
		case *typeinfo.Slot:
			if !f.Repeats {
				err = visitSlot(w.Walk(), evt)
			} else {
				err = visitSlots(w.Walk(), evt)
			}
		case *typeinfo.Num, *typeinfo.Str:
			if !f.Repeats {
				err = evt.Value(w)
			} else {
				err = visitValues(w.Walk(), evt)
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
	evt.Repeat(w)
	for it := w; it.Next(); {
		if e := visitFlow(it.Walk(), evt); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		evt.End(w)
	}
	return
}

func visitSlots(w Iter, evt Events) (err error) {
	evt.Repeat(w)
	for it := w; it.Next(); {
		if e := visitSlot(it.Walk(), evt); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		evt.End(w)
	}
	return
}

func visitValues(w Iter, evt Events) (err error) {
	evt.Repeat(w)
	for it := w; it.Next(); {
		if e := evt.Value(it); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		evt.End(w)
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
