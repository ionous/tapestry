package compact

import (
	"git.sr.ht/~ionous/iffy/export/jsn"
	"github.com/ionous/errutil"
)

// records a single value
type comValue struct {
	m    *CompactMarshaler
	name string // for debugging
	out  interface{}
}

func (d *comValue) named() string {
	return d.name
}
func (d *comValue) readData() interface{} {
	return d.out
}
func (d *comValue) writeData(v interface{}) {
	if d.out != nil {
		d.m.Warning(errutil.New("can only write data once"))
	} else {
		d.out = v
	}
}
func (d *comValue) MapValues(lede, kind string) {
	d.m.Error(errutil.New("cant map values"))
}
func (d *comValue) MapKey(sig, field string) {
	d.m.Error(errutil.New("key only valid in a map"))
}
func (d *comValue) MapLiteral(field string) {
	d.m.Error(errutil.New("literal only valid in a map"))
}
func (d *comValue) PickValues(kind, choice string) {
	d.m.Error(errutil.New("cant pick values"))
}
func (d *comValue) RepeatValues(hint int) {
	d.m.Error(errutil.New("cant repeat values"))
}
func (d *comValue) EndValues() {
	d.m.Error(errutil.New("cant end values"))
}

// WriteValue generically posts a primitive to the current state.
func (d *comValue) WriteValue(kind string, value interface{}) {
	// note: while the owner of this comValue memory is technically top state ...
	// in go, that owner type is inaccessible in the aggregated element
	// so... we need to the state machine for the outermost version of ourselves.
	d.m.writeData(value)
}

func (d *comValue) WriteChoice(kind string, val jsn.Enumeration) {
	var out string
	if str, ok := val.FindChoice(); !ok {
		out = val.String()
	} else {
		out = str
	}
	d.m.WriteValue(kind, out)
}
func (d *comValue) SetCursor(id string) {
	// only blocks have "at" cursor values... i think....
	d.m.Error(errutil.New("cant set cursor"))
}
func (d *comValue) Warning(e error) {
	// record an error but don't terminate
	d.m.err = errutil.Append(d.m.err, e)
}
func (d *comValue) Error(e error) {
	// record an error and terminate
	d.m.err = errutil.Append(d.m.err, e)
	d.m.stack = nil
	d.m.changeState(&comNull{})
}
