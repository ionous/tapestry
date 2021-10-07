package dout

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"github.com/ionous/errutil"
)

// Chart - marker so callers can see where a machine pointer came from.
type Chart struct{ *chart.Machine }

// NewDetailedMarshaler create an empty serializer to produce detailed script data.
func NewDetailedMarshaler() Chart {
	return Chart{chart.NewEncoder(newBlock)}
}

func makeEnum(val chart.EnumMarshaler) (ret string) {
	if k, v := val.GetEnum(); len(k) > 0 {
		ret = k
	} else {
		ret = v
	}
	return
}

// detailed data represents even primitive values as a map of {id,type,value}.
func newValue(m *chart.Machine, next *chart.StateMix) *chart.StateMix {
	next.OnValue = func(typeName string, pv interface{}) {
		if enum, ok := pv.(chart.EnumMarshaler); ok {
			pv = makeEnum(enum)
		}
		m.Commit(detValue{
			Id:    m.FlushCursor(),
			Type:  typeName,
			Value: pv,
		})
	}
	return next
}

// blocks handle beginning new flows, swaps, or repeats
// end ( and how they collect data ) gets left to the caller
func newBlock(m *chart.Machine) *chart.StateMix {
	next := chart.NewReportingState(m)
	next.OnMap = func(lede, typeName string) bool {
		m.PushState(newFlow(m, detMap{
			Id:     m.FlushCursor(),
			Type:   typeName,
			Fields: make(map[string]interface{}),
		}))
		return true
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
	return next
}

// flows create a set of key-values pairs
// the flow is closed ( written ) with a call to EndValues()
// every flow pushes a brand new machine
func newFlow(m *chart.Machine, vals detMap) *chart.StateMix {
	next := newBlock(m)
	next.OnKey = func(_, key string) bool {
		m.ChangeState(newKey(m, *next, key, vals))
		return true
	}
	next.OnLiteral = func(field string) bool {
		m.MapKey("", field) // loops back to OnKey
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
	next := newValue(m, &prev)
	next.OnCommit = func(v interface{}) {
		vals.Fields[key] = v // write our key, value pair
		m.ChangeState(&prev)
	}
	return next
}

// every slice pushes a brand new machine
func newSlice(m *chart.Machine, vals []interface{}) *chart.StateMix {
	next := newValue(m, newBlock(m))
	next.OnCommit = func(v interface{}) {
		vals = append(vals, v) // write a new value into the slice
	}
	next.OnEnd = func() {
		m.FinishState(vals)
	}
	return next
}

// every slice pushes a brand new machine
func newSwap(m *chart.Machine, choice string, vals detMap) *chart.StateMix {
	next := newValue(m, newBlock(m))
	next.OnCommit = func(v interface{}) {
		// write our choice and change into an error checking state
		vals.Fields = map[string]interface{}{
			choice: v,
		}
		m.ChangeState(chart.NewBlockResult(m, vals))
	}
	// fix? what should an uncommitted choice write?
	next.OnEnd = func() {
		m.FinishState(vals)
	}
	return next
}
