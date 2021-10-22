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
	OnBlock  func(jsn.BlockType) bool
	OnMap    func(string, string) bool
	OnKey    func(string, string) bool
	OnSlot   func(string, jsn.SlotBlock) bool
	OnPick   func(string, jsn.SwapBlock) bool
	OnRepeat func(string, jsn.SliceBlock) bool
	OnEnd    func()
	OnValue  func(string, interface{})
	OnWarn   func(error)
	OnCommit func(interface{})
}

// base state handles simple reporting.
func NewReportingState(m *Machine) *StateMix {
	return &StateMix{OnWarn: m.Error}
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

func (d *StateMix) MarshalBlock(b jsn.BlockType) (okay bool) {
	if call := d.OnBlock; call != nil {
		okay = call(b)
	} else {
		// backwards compatibility from when there were separate functions for each block type
		switch block := b.(type) {
		case jsn.FlowBlock:
			okay = d.MapValues(block.GetLede(), block.GetType())
		case jsn.SwapBlock:
			okay = d.PickValues(block.GetType(), block)
		case jsn.SliceBlock:
			okay = d.RepeatValues(block.GetType(), block)
		case jsn.SlotBlock:
			okay = d.SlotValues(block.GetType(), block)
		default:
			d.Warn(errutil.New("unexpected map type %T", b))
		}
	}
	return
}
func (d *StateMix) MapValues(lede, typeName string) (okay bool) {
	if call := d.OnMap; call != nil {
		okay = call(lede, typeName)
	} else {
		d.Warn(errutil.New("unexpected map", lede, typeName))
	}
	return
}
func (d *StateMix) MarshalKey(key, field string) (okay bool) {
	if call := d.OnKey; call != nil {
		okay = call(key, field)
	} else {
		d.Warn(errutil.New("unexpected key", key, field))
	}
	return
}
func (d *StateMix) SlotValues(typeName string, val jsn.SlotBlock) (okay bool) {
	if call := d.OnSlot; call != nil {
		okay = call(typeName, val)
	} else {
		d.Warn(errutil.New("unexpected pick", typeName, val))
	}
	return
}
func (d *StateMix) PickValues(typeName string, val jsn.SwapBlock) (okay bool) {
	if call := d.OnPick; call != nil {
		okay = call(typeName, val)
	} else {
		d.Warn(errutil.New("unexpected pick", typeName, val))
	}
	return
}
func (d *StateMix) RepeatValues(typeName string, val jsn.SliceBlock) (okay bool) {
	if call := d.OnRepeat; call != nil {
		okay = call(typeName, val)
	} else {
		d.Warn(errutil.New("unexpected repeat", typeName, val))
	}
	return
}
func (d *StateMix) EndBlock() {
	if call := d.OnEnd; call != nil {
		call()
	} else {
		d.Warn(errutil.New("unexpected end"))
	}
}
func (d *StateMix) MarshalValue(typeName string, pv interface{}) (okay bool) {
	if call := d.OnValue; call != nil {
		call(typeName, pv)
		okay = true
	} else {
		d.Warn(errutil.New("unexpected value", typeName, pv))
	}
	return
}
func (d *StateMix) Warn(err error) {
	if call := d.OnWarn; call != nil {
		call(err)
	} else {
		panic(err)
	}
}

// NOTE: doesn't warn if not implemented.
func (d *StateMix) Commit(v interface{}) {
	if call := d.OnCommit; call != nil {
		call(v)
	}
}
