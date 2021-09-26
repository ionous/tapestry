package compact

import (
	"git.sr.ht/~ionous/iffy/export/jsn"
	"github.com/ionous/errutil"
)

type compactMarshaler interface {
	jsn.Marshaler
	// substates write a fully completed value into us.
	commit(value interface{})
}

type comState struct {
	jsn.MarshalMix
	onCommit func(interface{})
}

func (d *comState) commit(v interface{}) {
	if call := d.onCommit; call != nil {
		call(v)
	} else {
		d.Error(errutil.New("cant commit", v))
	}
}

// base state handles simple reporting.
func newBase(m *CompactMarshaler, next *comState) *comState {
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
		m.changeState(&comState{MarshalMix: jsn.MarshalMix{
			// absorb all other errors
			// ( all other fns are empty,so they'll error and also be eaten )
			OnError: func(error) {},
		}})
	}
	return next
}

func newValue(m *CompactMarshaler, next *comState) *comState {
	next.OnValue = func(kind string, value interface{}) {
		m.commit(value)
	}
	return next
}

func newBlock(m *CompactMarshaler) *comState {
	next := newBase(m, new(comState))
	// starts a series of key-values pairs
	// the flow is closed ( written ) with a call to EndValues()
	next.OnMap = func(lede, _ string) {
		m.pushState(newFlow(m, newFlowData(lede)))
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	next.OnPick = func(kind, choice string) {
		m.pushState(newSwap(m))
	}
	next.OnRepeat = func(hint int) {
		m.pushState(newSlice(m, make([]interface{}, 0, hint)))
	}
	// in case nothing is written.
	next.OnEnd = func() {
		m.finishState(nil)
	}
	return next
}

func newFlow(m *CompactMarshaler, d *flowData) *comState {
	next := newBlock(m)
	next.OnKey = func(key, _ string) {
		m.changeState(newKey(m, *next, d, key))
	}
	next.OnLiteral = func(field string) {
		if len(d.values) > 0 {
			m.Error(errutil.New("unexpected literal after map key:value"))
		} else {
			m.changeState(newLit(m))
		}
	}
	// EndValues ends the current state and commits its data to the parent state.
	next.OnEnd = func() {
		m.finishState(d.finalize())
	}
	return next
}

// writes the value into the key and change back to the flow state
func newKey(m *CompactMarshaler, prev comState, d *flowData, key string) *comState {
	next := newValue(m, &prev)
	next.onCommit = func(v interface{}) {
		d.addMsg(key, v)
		m.changeState(&prev)
	}
	return next
}

// a literal is a block like value that results in a single value
// only writing the value or ending the block succeed.
func newLit(m *CompactMarshaler) *comState {
	next := newValue(m, newBlock(m))
	next.onCommit = func(v interface{}) {
		m.changeState(newBlockResult(m, v))
	}
	return next
}

// every slice pushes a brand new machine
func newSlice(m *CompactMarshaler, vals []interface{}) *comState {
	next := newValue(m, newBlock(m))
	// a new value is being added to our slice
	next.onCommit = func(v interface{}) {
		vals = append(vals, v)
	}
	// the slice is done, write it to our parent whomever that is.
	next.OnEnd = func() {
		m.finishState(vals)
	}
	return next
}

func newSwap(m *CompactMarshaler) *comState {
	next := newBlock(m)
	// we don't want to lose the *kind* of the choice
	// so we do this specially
	next.OnValue = func(kind string, value interface{}) {
		m.changeState(newBlockResult(m,
			map[string]interface{}{
				kind + ":": value,
			}))
	}
	// record the swap choice and move to an error detection state
	next.onCommit = func(v interface{}) {
		m.changeState(newBlockResult(m, v))
	}
	return next
}

// wait until the block is closed then finish
func newBlockResult(m *CompactMarshaler, v interface{}) *comState {
	return &comState{MarshalMix: jsn.MarshalMix{
		OnEnd: func() {
			m.finishState(v)
		},
	}}
}
