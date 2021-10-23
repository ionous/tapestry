package cout

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"github.com/ionous/errutil"
)

// xEncoder - marker so callers can see where a machine pointer came from.
type xEncoder struct {
	chart.Machine
}

// NewEncoder create an empty serializer to produce compact script data.
func Encode(in jsn.Marshalee) (ret interface{}, err error) {
	m := xEncoder{Machine: chart.MakeEncoder(custom)}
	next := m.newBlock()
	next.OnCommit = func(v interface{}) {
		if ret != nil {
			m.Error(errutil.New("can only write data once"))
		} else {
			ret = v
		}
	}
	m.ChangeState(next)
	in.Marshal(&m)
	return ret, m.Errors()
}

// debug.FactorialStory.Marshal(out)
// if d, e := out.Data(); e != nil {

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
func (m *xEncoder) newValue(next *chart.StateMix) *chart.StateMix {
	next.OnValue = func(typeName string, pv interface{}) error {
		m.Commit(unpack(pv))
		return nil
	}
	return next
}

// blocks handle beginning new flows, swaps, or repeats
// end ( and how they collect data ) gets left to the caller
func (m *xEncoder) newBlock() *chart.StateMix {
	var next chart.StateMix
	return m.addBlock(&next)
}

func (m *xEncoder) addBlock(next *chart.StateMix) *chart.StateMix {
	// starts a series of key-values pairs
	// the flow is closed ( written ) with a call to EndValues()
	next.OnMap = func(lede, _ string) bool {
		m.PushState(m.newFlow(newComFlow(lede)))
		return true
	}
	next.OnSlot = func(typeName string, slot jsn.Spotter) (okay bool) {
		if slot.HasSlot() {
			m.PushState(m.newSlot())
			okay = true
		}
		return
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	next.OnPick = func(typeName string, p jsn.Picker) (okay bool) {
		if choice, ok := p.GetChoice(); !ok {
			m.Error(errutil.New("couldnt determine choice of", p))
		} else if len(choice) > 0 {
			m.PushState(m.newSwap(&comSwap{typeName: typeName, choice: choice}))
			okay = true
		}
		return okay
	}
	next.OnRepeat = func(t string, vs jsn.Slicer) (okay bool) {
		if hint := vs.GetSize(); hint > 0 {
			m.PushState(m.newSlice(make([]interface{}, 0, hint)))
			okay = true
		}
		return okay
	}
	// next.OnEnd... gets determined by the specific block
	return next
}

func (m *xEncoder) newFlow(d *comFlow) *chart.StateMix {
	var next chart.StateMix
	next.OnKey = func(key, _ string) error {
		m.ChangeState(m.newKey(next, d, key))
		return nil
	}
	next.OnEnd = func() {
		// doesnt worry if there's a pending key/value
		// writing a value to a key is always considered optional
		m.FinishState(d.finalize())
	}
	return &next
}

// writes the value into the key and change back to the flow state
func (m *xEncoder) newKey(prev chart.StateMix, d *comFlow, key string) *chart.StateMix {
	next := m.newValue(m.addBlock(&prev))
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
func (m *xEncoder) newSlice(vals []interface{}) *chart.StateMix {
	next := m.newValue(m.newBlock())
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
func (m *xEncoder) newSlot() *chart.StateMix {
	next := m.newValue(m.newBlock())
	next.OnCommit = func(v interface{}) {
		// write our choice and change into an error checking state
		m.ChangeState(chart.NewBlockResult(&m.Machine, v))
	}
	// fix? what should an uncommitted choice write?
	next.OnEnd = func() {
		m.FinishState(nil)
	}
	return next
}

func (m *xEncoder) newSwap(swap *comSwap) *chart.StateMix {
	next := m.newValue(m.newBlock())
	next.OnCommit = func(v interface{}) {
		m.ChangeState(chart.NewBlockResult(&m.Machine, swap.SetValue(v)))
	}
	return next
}
