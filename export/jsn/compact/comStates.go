package compact

import (
	"git.sr.ht/~ionous/iffy/export/jsn/chart"
	"github.com/ionous/errutil"
)

// Chart - marker so callers can see where a machine pointer came from.
type Chart struct{ *chart.Machine }

// NewCompactMarshaler create an empty serializer to produce compact script data.
func NewCompactMarshaler() Chart {
	return Chart{chart.NewMachine(newBlock)}
}

// generically commits primitive value(s)
// compact data represents primitive values as their value.
func newValue(m *chart.Machine, next *chart.StateMix) *chart.StateMix {
	next.OnValue = func(kind string, value interface{}) {
		m.Commit(value)
	}
	return next
}

func newBlock(m *chart.Machine) *chart.StateMix {
	next := chart.NewReportingState(m)
	// starts a series of key-values pairs
	// the flow is closed ( written ) with a call to EndValues()
	next.OnMap = func(lede, _ string) {
		m.PushState(newFlow(m, newFlowData(lede)))
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	next.OnPick = func(kind, choice string) {
		m.PushState(newSwap(m))
	}
	next.OnRepeat = func(hint int) {
		m.PushState(newSlice(m, make([]interface{}, 0, hint)))
	}
	// in case nothing is written.
	next.OnEnd = func() {
		m.FinishState(nil)
	}
	return next
}

func newFlow(m *chart.Machine, d *flowData) *chart.StateMix {
	next := newBlock(m)
	next.OnKey = func(key, _ string) {
		m.ChangeState(newKey(m, *next, d, key))
	}
	next.OnLiteral = func(field string) {
		if len(d.values) > 0 {
			m.Error(errutil.New("unexpected literal after map key:value"))
		} else {
			m.ChangeState(newLit(m))
		}
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
	next := newBlock(m)
	// we don't want to lose the *kind* of the choice
	// so we do this specially
	next.OnValue = func(kind string, value interface{}) {
		m.ChangeState(chart.NewBlockResult(m,
			map[string]interface{}{
				kind + ":": value,
			}))
	}
	// record the swap choice and move to an error detection state
	next.OnCommit = func(v interface{}) {
		m.ChangeState(chart.NewBlockResult(m, v))
	}
	return next
}
