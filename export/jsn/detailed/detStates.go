package detailed

import (
	"git.sr.ht/~ionous/iffy/export/jsn"
	"github.com/ionous/errutil"
)

type detailedMarshaler interface {
	jsn.Marshaler
	// substates write a fully completed value into us.
	commit(interface{})
}

type detState struct {
	jsn.MarshalMix
	onCommit func(interface{})
}

func (d *detState) commit(v interface{}) {
	if call := d.onCommit; call != nil {
		call(v)
	} else {
		d.Error(errutil.New("cant commit", v))
	}
}

// base state handles simple reporting.
func newBase(m *DetailedMarshaler, next *detState) *detState {
	// for now, overwrite without error checking.
	next.OnCursor = func(id string) {
		m.cursor = id
	}
	// record an error but don't terminate
	next.OnWarn = func(e error) {
		m.err = errutil.Append(m.err, e)
	}
	// record an error and terminate all existing stats
	next.OnError = func(e error) {
		m.err = errutil.Append(m.err, e)
		m.stack = nil
		m.changeState(&detState{MarshalMix: jsn.MarshalMix{
			// absorb all other errors
			// ( all other fns are empty,so they'll error and also be eaten )
			OnError: func(error) {},
		}})
	}
	return next
}

// blocks handle beginning new flows, swaps, or repeats
// end ( and how they collect data ) gets left to the caller
func newBlock(m *DetailedMarshaler) *detState {
	next := newBase(m, new(detState))
	next.OnMap = func(lede, kind string) {
		m.pushState(newFlow(m, detMap{
			Id:     m.flushCursor(),
			Type:   kind,
			Fields: make(map[string]interface{}),
		}))
	}
	next.OnPick = func(kind, choice string) {
		m.pushState(newSwap(m, choice, detMap{
			Id:   m.flushCursor(),
			Type: kind,
		}))
	}
	next.OnRepeat = func(hint int) {
		m.pushState(newSlice(m, make([]interface{}, 0, hint)))
	}
	return next
}

// generically commits primitive value(s)
func newValue(m *DetailedMarshaler, next *detState) *detState {
	next.OnValue = func(kind string, value interface{}) {
		m.commit(detValue{
			Id:    m.flushCursor(),
			Type:  kind,
			Value: value,
		})
	}
	return next
}

// flows create a set of key-values pairs
// the flow is closed ( written ) with a call to EndValues()
// every flow pushes a brand new machine
func newFlow(m *DetailedMarshaler, vals detMap) *detState {
	next := newBlock(m)
	next.OnKey = func(_, key string) {
		m.changeState(newKey(m, *next, key, vals))
	}
	next.OnLiteral = func(field string) {
		m.MapKey("", field) // loops back to OnKey
	}
	next.OnEnd = func() {
		// doesnt worry if there's a pending key/value
		// writing a value to a key is always considered optional
		m.finishState(vals)
	}
	return next
}

// all keys are considered optional, so we do everything prev does with some extrs.
// keys wait until they have a value, then write their data into their parent's data;
// returning to the parent state.
func newKey(m *DetailedMarshaler, prev detState, key string, vals detMap) *detState {
	next := newValue(m, &prev)
	next.onCommit = func(v interface{}) {
		vals.Fields[key] = v // write our key, value pair
		m.changeState(&prev)
	}
	return next
}

// every slice pushes a brand new machine
func newSlice(m *DetailedMarshaler, vals []interface{}) *detState {
	next := newValue(m, newBlock(m))
	next.onCommit = func(v interface{}) {
		vals = append(vals, v) // write a new value into the slice
	}
	next.OnEnd = func() {
		m.finishState(vals)
	}
	return next
}

// every slice pushes a brand new machine
func newSwap(m *DetailedMarshaler, choice string, vals detMap) *detState {
	next := newValue(m, newBlock(m))
	next.onCommit = func(v interface{}) {
		// write our choice and change into an error checking state
		vals.Fields = map[string]interface{}{
			choice: v,
		}
		m.changeState(newBlockResult(m, vals))
	}
	// fix? what should an uncommitted choice write?
	next.OnEnd = func() {
		m.finishState(vals)
	}
	return next
}

// wait until the block is closed then finish
func newBlockResult(m *DetailedMarshaler, v interface{}) *detState {
	return &detState{MarshalMix: jsn.MarshalMix{
		OnEnd: func() {
			m.finishState(v)
		},
	}}
}
