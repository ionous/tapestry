package detailed

import (
	"strconv"

	"git.sr.ht/~ionous/iffy/export/jsn"
	"github.com/ionous/errutil"
)

type detBaseState struct {
	m    *DetailedMarshaler
	out  interface{}
	name string
}

type detFlowState struct {
	detBaseState // every flow pushes a brand new machine
	data         detMap
}

type detKeyState struct {
	*detFlowState // parent state
	key           string
	name          string
}

var namer int

func newName(base string) string {
	return base + "-" + strconv.Itoa(namer)
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

func (d *detBaseState) named() string {
	return d.name
}

func (d *detBaseState) readData() interface{} {
	return d.out
}

func (d *detBaseState) writeData(v interface{}) {
	if d.out != nil {
		d.m.Warning(errutil.New("base state already data"))
	} else {
		d.out = v
	}
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
		detBaseState: detBaseState{
			m:    d.m,
			name: newName("flow"),
		},
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

// WriteValue generically posts a primitive to the current state.
func (d *detBaseState) WriteValue(kind string, value interface{}) {
	// note: while the owner of this detBaseState memory is technically top state ...
	// in go, that owner type is inaccessible in the aggregated element
	// so... we need to the state machine for the outermost version of ourselves.
	d.m.writeData(detValue{
		Id:    d.m.flushCursor(),
		Type:  kind,
		Value: value,
	})
}

func (d *detBaseState) WriteChoice(kind string, val jsn.Enumeration) {
	d.m.WriteValue(kind, val.String())
}

func (d *detBaseState) PickValues(kind, choice string) {
	d.m.pushState(&detSwapState{
		detBaseState: detBaseState{
			m:    d.m,
			name: newName("swap"),
		},
		choice: choice,
		data: detMap{
			Id:   d.m.flushCursor(),
			Type: kind,
		},
	})
}

func (d *detBaseState) RepeatValues(hint int) {
	d.m.pushState(&detSliceState{
		detBaseState: detBaseState{
			m:    d.m,
			name: newName("slice"),
		},
		values: make([]interface{}, 0, hint),
	})
}

// EndValues ends the current state and writeDatas its data to the parent state.
func (d *detBaseState) EndValues() {
	was := d.m.popState()         // gets rid of us
	d.m.writeData(was.readData()) // write our accumulated data to the new parent
	// again, use "was" not "d" to get to the outermost version of ourself
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

func (d *detFlowState) readData() interface{} {
	return &d.data
}

func (d *detFlowState) writeData(v interface{}) {
	d.m.Error(errutil.New("missing key when writing to a flow"))
}

func (d *detKeyState) named() string {
	return "key-" + d.detFlowState.named()
}

func (d *detKeyState) writeData(v interface{}) {
	d.data.Fields[d.key] = v // write our key, value pair
	d.m.changeState(d.detFlowState)
}

func (d *detSliceState) readData() interface{} {
	return d.values
}

func (d *detSliceState) writeData(v interface{}) {
	d.values = append(d.values, v)
}

func (d *detSwapState) writeData(v interface{}) {
	d.data.Fields = map[string]interface{}{
		d.choice: v,
	}
	d.m.changeState(&detSwapWritten{d})
}

func (d *detSwapState) readData() interface{} {
	return &d.data
}

type detSwapWritten struct {
	*detSwapState
}

func (d *detSwapWritten) writeData(v interface{}) {
	d.m.Warning(errutil.New("swap already written"))
}
