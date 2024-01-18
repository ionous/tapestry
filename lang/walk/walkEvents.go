package walk

import "errors"

// fix? this should probably be passing type info and not walker itself
type Events interface {
	Flow(Walker) error
	Slot(Walker) error
	Repeat(Walker) error
	// called for every member of a flow
	Field(Walker) error
	// called for each value in a flow that exists
	Value(Walker) error // FIX
	// called after each flow, slot, or repeat
	End(Walker) error
}

// can be used to early exit from visiting
// caller still has to check for this to disambiguate the returned error
var DoneVisiting = errors.New("done visiting")

// turn the iterator into an event style traversal
// ( helps especially with backwards compatibility )
func VisitFlow(w Walker, evt Events) error {
	return visitFlow(w, evt)
	// fix: if we looked at the outer type instead of the fields
	// this would be a lot more flexible
	// currently, requires a block
}

func VisitSlot(w Walker, evt Events) error {
	return visitSlot(w, evt)
}

func visitFields(w Walker, evt Events) (err error) {
	for w.Next() && err == nil {
		if e := evt.Field(w); e != nil {
			err = e
			break
		}
		switch f := w.Field(); f.SpecType() {
		case Flow:
			if !f.Repeats() {
				err = visitFlow(w.Walk(), evt)
			} else if e := visitFlows(w.Walk(), evt); e != nil {
				err = e
			} else {
				evt.End(w)
			}
		case Slot:
			if !f.Repeats() {
				err = visitSlot(w.Walk(), evt)
			} else if e := visitSlots(w.Walk(), evt); e != nil {
				err = e
			} else {
				evt.End(w)
			}
		case Value:
			if !f.Repeats() {
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

func visitFlow(w Walker, evt Events) (err error) {
	if e := evt.Flow(w); e != nil {
		err = e
	} else if e := visitFields(w, evt); e != nil {
		err = e
	} else {
		err = evt.End(w)
	}
	return
}

func visitSlot(w Walker, evt Events) (err error) {
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

func visitFlows(w Walker, evt Events) (err error) {
	for w.Next() {
		if e := visitFlow(w.Walk(), evt); e != nil {
			err = e
			break
		}
	}
	return
}

func visitSlots(w Walker, evt Events) (err error) {
	for w.Next() {
		if e := visitSlot(w.Walk(), evt); e != nil {
			err = e
			break
		}
	}
	return
}

func visitValues(w Walker, evt Events) (err error) {
	for w.Next() {
		if e := evt.Value(w); e != nil {
			err = e
			break
		}
	}
	return
}

// given a slot, return its flow ( false if the slot was empty )
func unpackSlot(slot Walker) (ret Walker, okay bool) {
	if slot.Next() { // next moves the focus to the flow if any
		ret, okay = slot.Walk(), true // walk returns the flow as a container
	}
	return
}
