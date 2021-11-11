package chart

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"github.com/ionous/errutil"
)

type State interface {
	jsn.State
	// substates write a fully completed value into us.
	// fix -- maybe commit should return error? not sure.
	Commit(interface{})
}

// Unhandled - an error code to break out of loops
type Unhandled string

func (e Unhandled) Error() string {
	return "Unhandled state during " + string(e)
}

// StateMix implements the jsn.State interface
// providing functions which can be overridden one at a time to customize functionality
type StateMix struct {
	OnBlock  func(jsn.Block) error
	OnMap    func(string, jsn.FlowBlock) bool
	OnKey    func(string, string) error
	OnSlot   func(string, jsn.SlotBlock) bool
	OnSwap   func(string, jsn.SwapBlock) bool
	OnRepeat func(string, jsn.SliceBlock) bool
	OnEnd    func()
	OnValue  func(string, interface{}) error
	OnCommit func(interface{})
}

// wait until the block is closed then finish
func NewBlockResult(m *Machine, v interface{}) *StateMix {
	return &StateMix{
		OnEnd: func() {
			m.FinishState(v)
		},
		OnCommit: func(interface{}) {
			m.Error(Unhandled("OnCommit"))
		},
	}
}

func (d *StateMix) MarshalBlock(b jsn.Block) (err error) {
	if call := d.OnBlock; call != nil {
		err = call(b)
	} else {
		// backwards compatibility from when there were separate functions for each block type
		switch block := b.(type) {
		case jsn.FlowBlock:
			err = d.MapValues(block.GetType(), block)
		case jsn.SwapBlock:
			err = d.PickValues(block.GetType(), block)
		case jsn.SliceBlock:
			err = d.RepeatValues(block.GetType(), block)
		case jsn.SlotBlock:
			err = d.SlotValues(block.GetType(), block)
		default:
			err = errutil.Fmt("unexpected map type %T", b)
		}
	}
	return
}
func okayMissing(ok bool) (err error) {
	if !ok {
		err = jsn.Missing
	}
	return
}
func (d *StateMix) MapValues(typeName string, val jsn.FlowBlock) (err error) {
	if call := d.OnMap; call != nil {
		err = okayMissing(call(typeName, val))
	} else {
		err = Unhandled(errutil.Sprint("MapValues", typeName, val))
	}
	return
}
func (d *StateMix) MarshalKey(key, field string) (err error) {
	if call := d.OnKey; call != nil {
		err = call(key, field)
	} else {
		err = Unhandled(errutil.Sprint("MarshalKey", key, field))
	}
	return
}
func (d *StateMix) SlotValues(typeName string, val jsn.SlotBlock) (err error) {
	if call := d.OnSlot; call != nil {
		err = okayMissing(call(typeName, val))
	} else {
		err = Unhandled(errutil.Sprint("SlotValues", typeName, val))
	}
	return
}
func (d *StateMix) PickValues(typeName string, val jsn.SwapBlock) (err error) {
	if call := d.OnSwap; call != nil {
		err = okayMissing(call(typeName, val))
	} else {
		err = Unhandled(errutil.Sprint("PickValues", typeName, val))
	}
	return
}
func (d *StateMix) RepeatValues(typeName string, val jsn.SliceBlock) (err error) {
	if call := d.OnRepeat; call != nil {
		err = okayMissing(call(typeName, val))
	} else {
		err = Unhandled(errutil.Sprint("RepeatValues", typeName, val))
	}
	return
}
func (d *StateMix) EndBlock() {
	if call := d.OnEnd; call != nil {
		call()
	} /* else {
		// no way to report on this any more :/
		err = Unhandled(errutil.Sprint("EndBlock"))
	}*/
}
func (d *StateMix) MarshalValue(typeName string, pv interface{}) (err error) {
	if call := d.OnValue; call != nil {
		err = call(typeName, pv)
	} else {
		err = Unhandled(errutil.Sprint("MarshalValue", typeName, pv))
	}
	return
}

// NOTE: doesn't warn if not implemented.
func (d *StateMix) Commit(v interface{}) {
	if call := d.OnCommit; call != nil {
		call(v)
	}
}
