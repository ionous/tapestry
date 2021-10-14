package chart

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"github.com/ionous/errutil"
)

type State interface {
	jsn.Marshaler
	// substates write a fully completed value into us.
	Commit(interface{})
}

type StateMix struct {
	jsn.MarshalMix
	OnCommit func(interface{})
}

func (d *StateMix) Commit(v interface{}) {
	if call := d.OnCommit; call != nil {
		call(v)
	} else {
		d.Error(errutil.New("cant commit", v))
	}
}

// base state handles simple reporting.
func NewReportingState(m *Machine) *StateMix {
	next := new(StateMix)
	// for now, overwrite without error checking.
	next.OnCursor = func(id string) {
		m.cursor = id
	}
	// record an error but don't terminate
	next.OnWarn = func(e error) {
		m.err = errutil.Append(m.err, e)
	}
	return next
}

// wait until the block is closed then finish
func NewBlockResult(m *Machine, v interface{}) *StateMix {
	return &StateMix{MarshalMix: jsn.MarshalMix{
		OnEnd: func() {
			m.FinishState(v)
		},
	}}
}
