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

// StateMix implements the jsn.State interface
// providing functions which can be overridden one at a time to customize functionality
type StateMix struct {
	OnBlock  func(jsn.BlockType) error
	OnMap    func(string, string) bool
	OnKey    func(string, string) bool
	OnSlot   func(string, jsn.SlotBlock) bool
	OnPick   func(string, jsn.SwapBlock) bool
	OnRepeat func(string, jsn.SliceBlock) bool
	OnEnd    func()
	OnValue  func(string, interface{})
	OnCommit func(interface{})
}

// base state handles simple reporting.
func NewReportingState(m *Machine) *StateMix {
	return &StateMix{}
}

// wait until the block is closed then finish
func NewBlockResult(m *Machine, v interface{}) *StateMix {
	return &StateMix{
		OnEnd: func() {
			m.FinishState(v)
		},
		OnCommit: func(interface{}) {
			m.Error(errutil.New("expected a terminal value"))
		},
	}
}

func (d *StateMix) MarshalBlock(b jsn.BlockType) (err error) {
	if call := d.OnBlock; call != nil {
		err = call(b)
	} else {
		// backwards compatibility from when there were separate functions for each block type
		switch block := b.(type) {
		case jsn.FlowBlock:
			err = d.MapValues(block.GetLede(), block.GetType())
		case jsn.SwapBlock:
			err = d.PickValues(block.GetType(), block)
		case jsn.SliceBlock:
			err = d.RepeatValues(block.GetType(), block)
		case jsn.SlotBlock:
			err = d.SlotValues(block.GetType(), block)
		default:
			err = errutil.New("unexpected map type %T", b)
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
func (d *StateMix) MapValues(lede, typeName string) (err error) {
	if call := d.OnMap; call != nil {
		err = okayMissing(call(lede, typeName))
	} else {
		err = errutil.New("unexpected map", lede, typeName)
	}
	return
}
func (d *StateMix) MarshalKey(key, field string) (err error) {
	if call := d.OnKey; call != nil {
		err = okayMissing(call(key, field))
	} else {
		err = errutil.New("unexpected key", key, field)
	}
	return
}
func (d *StateMix) SlotValues(typeName string, val jsn.SlotBlock) (err error) {
	if call := d.OnSlot; call != nil {
		err = okayMissing(call(typeName, val))
	} else {
		err = errutil.New("unexpected pick", typeName, val)
	}
	return
}
func (d *StateMix) PickValues(typeName string, val jsn.SwapBlock) (err error) {
	if call := d.OnPick; call != nil {
		err = okayMissing(call(typeName, val))
	} else {
		err = errutil.New("unexpected pick", typeName, val)
	}
	return
}
func (d *StateMix) RepeatValues(typeName string, val jsn.SliceBlock) (err error) {
	if call := d.OnRepeat; call != nil {
		err = okayMissing(call(typeName, val))
	} else {
		err = errutil.New("unexpected repeat", typeName, val)
	}
	return
}
func (d *StateMix) EndBlock() {
	if call := d.OnEnd; call != nil {
		call()
	} /* else {
		// no way to report on this any more :/
		err = errutil.New("unexpected end"))
	}*/
}
func (d *StateMix) MarshalValue(typeName string, pv interface{}) (err error) {
	if call := d.OnValue; call != nil {
		call(typeName, pv)
	} else {
		err = errutil.New("unexpected value", typeName, pv)
	}
	return
}

// NOTE: doesn't warn if not implemented.
func (d *StateMix) Commit(v interface{}) {
	if call := d.OnCommit; call != nil {
		call(v)
	}
}
