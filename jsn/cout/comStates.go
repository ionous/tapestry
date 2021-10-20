package cout

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"github.com/ionous/errutil"
)

// Encoder - marker so callers can see where a machine pointer came from.
type Encoder struct{ *chart.Machine }

// NewEncoder create an empty serializer to produce compact script data.
func NewEncoder(cs chart.Customization) *Encoder {
	return &Encoder{chart.NewEncoder(cs, newBlock)}
}

func unpack(pv interface{}) (ret interface{}) {
	switch pv := pv.(type) {
	case interface{ GetCompactValue() interface{} }:
		ret = pv.GetCompactValue()
	case interface{ GetValue() interface{} }:
		ret = pv.GetValue()
	default:
		ret = pv
	}
	return
}

// compact data represents primitive values as their value.
func newValue(m *chart.Machine, next *chart.StateMix) *chart.StateMix {
	next.OnValue = func(typeName string, pv interface{}) {
		m.Commit(unpack(pv))
	}
	return next
}

// blocks handle beginning new flows, swaps, or repeats
// end ( and how they collect data ) gets left to the caller
func newBlock(m *chart.Machine) *chart.StateMix {
	return addBlock(m, chart.NewReportingState(m))
}

func addBlock(m *chart.Machine, next *chart.StateMix) *chart.StateMix {
	// starts a series of key-values pairs
	// the flow is closed ( written ) with a call to EndValues()
	next.OnMap = func(lede, _ string) bool {
		m.PushState(newFlow(m, newComFlow(lede, false)))
		return true
	}
	next.OnLiteral = func(lede, _ string) bool {
		m.PushState(newFlow(m, newComFlow(lede, true)))
		return true
	}
	next.OnSlot = func(typeName string, slot jsn.Spotter) (okay bool) {
		if slot.HasSlot() {
			m.PushState(newSlot(m))
			okay = true
		}
		return
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	next.OnPick = func(typeName string, p jsn.Picker) (okay bool) {
		if choice, ok := p.GetChoice(); !ok {
			m.Error(errutil.New("couldnt determine choice of", p))
		} else if len(choice) > 0 {
			m.PushState(newSwap(m, &comSwap{typeName: typeName, choice: choice}))
			okay = true
		}
		return okay
	}
	next.OnRepeat = func(t string, vs jsn.Slicer) (okay bool) {
		if hint := vs.GetSize(); hint > 0 {
			m.PushState(newSlice(m, make([]interface{}, 0, hint)))
			okay = true
		}
		return okay
	}
	// next.OnEnd... gets determined by the specific block
	return next
}

func newFlow(m *chart.Machine, d *comFlow) *chart.StateMix {
	next := chart.NewReportingState(m)
	next.OnKey = func(key, _ string) bool {
		m.ChangeState(newKey(m, *next, d, key))
		return true
	}
	next.OnEnd = func() {
		// doesnt worry if there's a pending key/value
		// writing a value to a key is always considered optional
		m.FinishState(d.finalize())
	}
	return next
}

// writes the value into the key and change back to the flow state
func newKey(m *chart.Machine, prev chart.StateMix, d *comFlow, key string) *chart.StateMix {
	next := newValue(m, addBlock(m, &prev))
	next.OnCommit = func(v interface{}) {
		if c, ok := v.(*comSwap); ok {
			d.addMsgPair(key, c.choice, c.value)
		} else {
			d.addMsg(key, v)
			m.ChangeState(&prev)
		}
	}
	return next
}

// every slice pushes a brand new machine
func newSlice(m *chart.Machine, vals []interface{}) *chart.StateMix {
	next := newValue(m, newBlock(m))
	// a new value is being added to our slice
	next.OnCommit = func(v interface{}) {
		vals = append(vals, v)
	}
	// the slice is done, write it to our parent whomever that is.
	next.OnEnd = func() {
		m.FinishState(vals)
	}
	return next
}

// every slot pushes a brand new machine
func newSlot(m *chart.Machine) *chart.StateMix {
	next := newValue(m, newBlock(m))
	next.OnCommit = func(v interface{}) {
		// write our choice and change into an error checking state
		m.ChangeState(chart.NewBlockResult(m, v))
	}
	// fix? what should an uncommitted choice write?
	next.OnEnd = func() {
		m.FinishState(nil)
	}
	return next
}

func newSwap(m *chart.Machine, swap *comSwap) *chart.StateMix {
	next := newValue(m, newBlock(m))
	next.OnCommit = func(v interface{}) {
		m.ChangeState(chart.NewBlockResult(m, swap.SetValue(v)))
	}
	return next
}
