package compact

import (
	"git.sr.ht/~ionous/iffy/export/jsn"
	"git.sr.ht/~ionous/iffy/export/jsn/chart"
	"github.com/ionous/errutil"
)

// Chart - marker so callers can see where a machine pointer came from.
type Chart struct{ *chart.Machine }

// NewCompactMarshaler create an empty serializer to produce compact script data.
func NewCompactMarshaler() Chart {
	return Chart{chart.NewMachine(newBlock)}
}

// compact data represents primitive values as their value.
func newValue(m *chart.Machine, next *chart.StateMix) *chart.StateMix {
	return chart.OnValue(next, func(k string, v interface{}) {
		m.Commit(v)
	})
}

func newBlock(m *chart.Machine) *chart.StateMix {
	next := chart.NewReportingState(m)
	// starts a series of key-values pairs
	// the flow is closed ( written ) with a call to EndValues()
	next.OnMap = func(lede, _ string) bool {
		m.PushState(newFlow(m, newFlowData(lede)))
		return true
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	// the compact encoding relies on the encoded inner block to unpack the choice.
	// ( implies each option needs to be a unique type. )
	next.OnPick = func(p jsn.Picker) (okay bool) {
		if c, ok := p.GetChoice(); !ok {
			m.Error(errutil.New("couldnt determine choice of", p))
		} else if len(c) > 0 {
			m.PushState(newSwap(m))
			okay = true
		}
		return okay
	}
	next.OnRepeat = func(hint int) (okay bool) {
		if hint > 0 {
			m.PushState(newSlice(m, make([]interface{}, 0, hint)))
			okay = true
		}
		return okay
	}
	// in case nothing is written.
	next.OnEnd = func() {
		m.FinishState(nil)
	}
	return next
}

func newFlow(m *chart.Machine, d *flowData) *chart.StateMix {
	next := newBlock(m)
	next.OnKey = func(key, _ string) bool {
		m.ChangeState(newKey(m, *next, d, key))
		return true
	}
	next.OnLiteral = func(field string) bool {
		if len(d.values) > 0 {
			m.Error(errutil.New("unexpected literal after map key:value"))
		} else {
			m.ChangeState(newLit(m))
		}
		return true
	}
	// EndValues ends the current state and commits its data to the parent state.
	next.OnEnd = func() {
		m.FinishState(d.finalize())
	}
	return next
}

// writes the value into the key and change back to the flow state
func newKey(m *chart.Machine, prev chart.StateMix, d *flowData, key string) *chart.StateMix {
	next := newValue(m, &prev)
	next.OnCommit = func(v interface{}) {
		d.addMsg(key, v)
		m.ChangeState(&prev)
	}
	return next
}

// a literal is a block like value that results in a single value
// only writing the value or ending the block succeed.
func newLit(m *chart.Machine) *chart.StateMix {
	next := newValue(m, newBlock(m))
	next.OnCommit = func(v interface{}) {
		m.ChangeState(chart.NewBlockResult(m, v))
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

func newSwap(m *chart.Machine) *chart.StateMix {
	// we don't want to lose the *kind* of the choice for simple values
	next := chart.OnValue(newBlock(m), func(kind string, value interface{}) {
		m.ChangeState(chart.NewBlockResult(m,
			map[string]interface{}{
				kind + ":": value,
			}))
	})
	// record the swap choice and move to an error detection state
	next.OnCommit = func(v interface{}) {
		m.ChangeState(chart.NewBlockResult(m, v))
	}
	return next
}
