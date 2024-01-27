package inspect

import (
	"errors"
	"fmt"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// fix? this should probably be passing type info and not walker itself
type Events interface {
	Flow(It) error
	// called for every member of a flow.
	Field(It) error
	// called for a member slot: can be an empty slot.
	Slot(It) error
	// called for a member that repeats; can be an empty list.
	Repeat(It) error
	// called for each str or num in a flow, or in a repeat.
	// because it isnt a container, there is no matching end.
	Value(It) error
	// called after each flow, slot, or repeat.
	End(It) error
}

// can be used to early exit from visiting
// caller still has to check for this to disambiguate the returned error
var DoneVisiting = errors.New("done visiting")

// turn the iterator into an event style callbacks
func Visit(i typeinfo.Instance, evt Events) (err error) {
	w := Walk(i) // cheats a little to grab the type info and repeat style...
	t, repeat := w.currType, w.curr.Kind() == r.Slice
	// almost the same as visit fields except it has to walk into the field
	// and this already has the element it walks to walk
	switch t.(type) {
	case *typeinfo.Flow:
		if !repeat {
			err = visitFlow(w, evt)
		} else {
			err = visitFlows(w, evt)
		}
	case *typeinfo.Slot:
		if !repeat {
			err = visitSlot(w, evt)
		} else {
			err = visitSlots(w, evt)
		}
	default:
		err = fmt.Errorf("expected a container type, not %T", i)
	}
	return
}

func visitFields(w It, evt Events) (err error) {
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

func visitFlow(w It, evt Events) (err error) {
	if e := evt.Flow(w); e != nil {
		err = e
	} else if e := visitFields(w, evt); e != nil {
		err = e
	} else {
		err = evt.End(w)
	}
	return
}

func visitSlot(w It, evt Events) (err error) {
	if e := evt.Slot(w); e != nil {
		err = e
	} else {
		if it := w; it.Next() { // next moves the focus to the flow if any
			err = visitFlow(it.Walk(), evt) // walk returns the flow as a container
		}
		if err == nil {
			err = evt.End(w)
		}
	}
	return
}

func visitFlows(w It, evt Events) (err error) {
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

func visitSlots(w It, evt Events) (err error) {
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

func visitValues(w It, evt Events) (err error) {
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
