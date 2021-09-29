package chart

import (
	"git.sr.ht/~ionous/iffy/export/jsn"
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
	// record an error and terminate all existing stats
	next.OnError = func(e error) {
		m.err = errutil.Append(m.err, e)
		m.stack = nil
		m.ChangeState(&StateMix{MarshalMix: jsn.MarshalMix{
			// absorb all other errors
			// ( all other fns are empty,so they'll error and also be eaten )
			OnError: func(error) {},
		}})
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

// OnValue generically handles primitive value(s)
func OnValue(next *StateMix, onValue func(string, interface{})) *StateMix {
	next.OnBool = func(val jsn.BoolMarshaler) {
		onValue(val.GetType(), val.GetBool())
	}
	next.OnEnum = func(val jsn.EnumMarshaler) {
		var out string
		if k, v := val.GetEnum(); len(k) > 0 {
			out = k
		} else {
			out = v
		}
		onValue(val.GetType(), out)
	}
	next.OnNum = func(val jsn.NumMarshaler) {
		onValue(val.GetType(), val.GetNum())
	}
	next.OnStr = func(val jsn.StrMarshaler) {
		onValue(val.GetType(), val.GetStr())
	}
	return next
}
