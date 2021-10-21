package chart

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"github.com/ionous/errutil"
)

type State interface {
	jsn.State
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
		d.Warn(errutil.Fmt("unexpected commit (%v)", v))
	}
}

// base state handles simple reporting.
func NewReportingState(m *Machine) *StateMix {
	return &StateMix{MarshalMix: jsn.MarshalMix{OnWarn: m.Error}}
}

// wait until the block is closed then finish
func NewBlockResult(m *Machine, v interface{}) *StateMix {
	return &StateMix{MarshalMix: jsn.MarshalMix{
		OnEnd: func() {
			m.FinishState(v)
		},
	}}
}
