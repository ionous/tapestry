package dout

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"github.com/ionous/errutil"
)

// Chart - marker so callers can see where a machine pointer came from.
type xEncoder struct{ chart.Machine }

// NewEncoder create an empty serializer to produce detailed script data.
func Encode(in jsn.Marshalee) (ret interface{}, err error) {
	m := xEncoder{chart.MakeEncoder(nil)}
	next := newBlock(&m.Machine)
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

func newValue(m *chart.Machine, next *chart.StateMix) *chart.StateMix {
	next.OnValue = func(typeName string, pv interface{}) {
		if el, ok := pv.(interface{ GetValue() interface{} }); ok {
			pv = el.GetValue()
		}
		// detailed data represents even primitive values as a map of {id,type,value}.
		m.Commit(detValue{
			Id:    m.FlushCursor(),
			Type:  typeName,
			Value: pv,
		})
	}
	// next.OnCommit -- handled by each caller
	return next
}

// blocks handle beginning new flows, swaps, or repeats
// end ( and how they collect data ) gets left to the caller
func newBlock(m *chart.Machine) *chart.StateMix {
	return addBlock(m, chart.NewReportingState(m))
}

func addBlock(m *chart.Machine, next *chart.StateMix) *chart.StateMix {
	next.OnMap = func(_, typeName string) bool {
		m.PushState(newFlow(m, detMap{
			Id:     m.FlushCursor(),
			Type:   typeName,
			Fields: make(map[string]interface{}),
		}))
		return true
	}
	next.OnSlot = func(typeName string, slot jsn.Spotter) (okay bool) {
		if slot.HasSlot() {
			m.PushState(newSlot(m, detValue{
				Id:   m.FlushCursor(),
				Type: typeName,
			}))
			okay = true
		}
		return
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	next.OnPick = func(typeName string, p jsn.Picker) (okay bool) {
		if choice, ok := p.GetChoice(); !ok {
			m.Error(errutil.New("couldnt determine choice of", p))
		} else if len(choice) > 0 {
			m.PushState(newSwap(m, choice, detMap{
				Id:   m.FlushCursor(),
				Type: typeName,
			}))
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

// flows create a set of key-values pairs
// the flow is closed ( written ) with a call to EndValues()
// every flow pushes a brand new machine
func newFlow(m *chart.Machine, vals detMap) *chart.StateMix {
	next := chart.NewReportingState(m)
	next.OnKey = func(_, key string) bool {
		m.ChangeState(newKey(m, *next, key, vals))
		return true
	}
	next.OnEnd = func() {
		// doesnt worry if there's a pending key/value
		// writing a value to a key is always considered optional
		m.FinishState(vals)
	}
	return next
}

// all keys are considered optional, so we do everything prev does with some extrs.
// keys wait until they have a value, then write their data into their parent's data;
// returning to the parent state.
func newKey(m *chart.Machine, prev chart.StateMix, key string, vals detMap) *chart.StateMix {
	// a key's value can be a simple value, or a block.
	next := newValue(m, addBlock(m, &prev))
	next.OnCommit = func(v interface{}) {
		vals.Fields[key] = v // write our key, value pair
		m.ChangeState(&prev)
	}
	return next
}

// every slice pushes a brand new machine
func newSlice(m *chart.Machine, slice []interface{}) *chart.StateMix {
	next := newValue(m, newBlock(m))
	next.OnCommit = func(v interface{}) {
		slice = append(slice, v) // write a new value into the slice
	}
	next.OnEnd = func() {
		m.FinishState(slice)
	}
	return next
}

// every swap pushes a brand new machine
func newSlot(m *chart.Machine, slot detValue) *chart.StateMix {
	next := newValue(m, newBlock(m))
	next.OnCommit = func(v interface{}) {
		// write our choice and change into an error checking state
		slot.Value = v
		m.ChangeState(chart.NewBlockResult(m, slot))
	}
	// fix? what should an uncommitted choice write?
	next.OnEnd = func() {
		m.FinishState(slot)
	}
	return next
}

// every swap pushes a brand new machine
func newSwap(m *chart.Machine, choice string, swap detMap) *chart.StateMix {
	next := newValue(m, newBlock(m))
	next.OnCommit = func(v interface{}) {
		// write our choice and change into an error checking state
		swap.Fields = map[string]interface{}{
			choice: v,
		}
		m.ChangeState(chart.NewBlockResult(m, swap))
	}
	// fix? what should an uncommitted choice write?
	next.OnEnd = func() {
		m.FinishState(swap)
	}
	return next
}
