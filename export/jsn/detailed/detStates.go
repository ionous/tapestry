package detailed

import (
	"git.sr.ht/~ionous/iffy/export/jsn"
	"github.com/ionous/errutil"
)

type detBaseState struct {
	m *DetailedMarshaler
}

type detFlowState struct {
	detBaseState // every flow pushes a brand new machine
	data         detMap
}

type detKeyState struct {
	*detFlowState // parent state
	key           string
}

type detSliceState struct {
	detBaseState // every slice pushes a brand new machine
	values       []interface{}
}

type detSwapState struct {
	detBaseState // every swap pushes a brand new machine
	data         detMap
	choice       string
}

func (d *detBaseState) EndValues() {
	panic("end values not implemented")
}

// record an error but don't terminate
func (d *detBaseState) Warning(e error) {
	d.m.err = errutil.Append(d.m.err, e)
}

// record an error and terminate
func (d *detBaseState) Error(e error) {
	d.m.err = errutil.Append(d.m.err, e)
	d.m.stack = nil
	d.m.changeState(detNull{})
}

// starts a series of key-values pairs
// the flow is closed ( written ) with a call to EndValues()
func (d *detBaseState) MapValues(lede, kind string) {
	d.m.pushState(&detFlowState{
		detBaseState: detBaseState{m: d.m},
		data: detMap{
			Id:     d.m.flushCursor(),
			Type:   kind,
			Fields: make(map[string]interface{}),
		},
	})
}

func (d *detBaseState) MapKey(sig, field string) {
	d.m.Error(errutil.New("key only valid in a map"))
}

func (d *detBaseState) MapLiteral(string) {
	d.m.Error(errutil.New("literal only valid in a map"))
}

func (d *detBaseState) SetCursor(id string) {
	d.m.cursor = id //  for now, overwrite without error checking.
}

// SpecifyValue generically posts a primitive to the current state.
func (d *detBaseState) SpecifyValue(kind string, value interface{}) {
	// note: while the owner of this detBaseState memory is technically top state ...
	// in go, that owner type is inaccessible in the aggregated element
	// so... we need to the state machine for the outermost version of ourselves.
	d.m.commit(detValue{
		Id:    d.m.flushCursor(),
		Type:  kind,
		Value: value,
	})
}

func (d *detBaseState) SpecifyEnum(kind string, val jsn.Enumeration) {
	d.m.SpecifyValue(kind, val.String())
}

func (d *detBaseState) PickValues(kind, choice string) {
	d.m.pushState(&detSwapState{
		detBaseState: detBaseState{m: d.m},
		choice:       choice,
		data: detMap{
			Id:   d.m.flushCursor(),
			Type: kind,
		},
	})
}

func (d *detBaseState) RepeatValues(hint int) {
	d.m.pushState(&detSliceState{
		detBaseState: detBaseState{m: d.m},
		values:       make([]interface{}, 0, hint),
	})
}

func (d *detFlowState) MapKey(sig, field string) {
	d.m.changeState(&detKeyState{
		detFlowState: d,
		key:          field,
	})
}

func (d *detFlowState) MapLiteral(field string) {
	d.m.MapKey("", field)
}

// EndValues ends the current state and commits its data to the parent state.
func (d *detFlowState) EndValues() {
	d.m.finishState(d.data)
}

func (d *detFlowState) commit(v interface{}) {
	d.m.Error(errutil.New("missing key when writing to a flow"))
}

func (d *detKeyState) commit(v interface{}) {
	d.data.Fields[d.key] = v // write our key, value pair
	d.m.changeState(d.detFlowState)
}

// write a new value into the slice
func (d *detSliceState) commit(v interface{}) {
	d.values = append(d.values, v)
}

// EndValues ends the current state and commits its data to the parent state.
func (d *detSliceState) EndValues() {
	d.m.finishState(d.values)
}

// write our choice and change into an error checking state
func (d *detSwapState) commit(v interface{}) {
	d.data.Fields = map[string]interface{}{
		d.choice: v,
	}
	d.m.changeState(&detSwapWritten{d})
}

// EndValues ends the current state and commits its data to the parent state.
func (d *detSwapState) EndValues() {
	d.m.finishState(d.data)
}

type detSwapWritten struct {
	*detSwapState
}

func (d *detSwapWritten) commit(v interface{}) {
	d.m.Warning(errutil.New("swap already committed"))
}
