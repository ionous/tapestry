package compact

import (
	"github.com/ionous/errutil"
)

type comFlow struct {
	comBlock // every flow pushes a brand new machine
	sig      Sig
	values   []interface{}
}

type comKey struct {
	*comFlow // parent state
	key      string
}

type comLiteral struct {
	comValue
}

type comSlice struct {
	comBlock // every slice pushes a brand new machine
	values   []interface{}
}

type comSwap struct {
	comBlock // every swap pushes a brand new machine
}

func (d *comFlow) MapKey(sig, _ string) {
	d.m.changeState(&comKey{
		comFlow: d,
		key:     sig,
	})
}

func (d *comFlow) MapLiteral(field string) {
	if len(d.values) > 0 {
		d.m.Error(errutil.New("unexpected literal after map key:value"))
	} else {
		d.m.changeState(&comLiteral{comValue: comValue{m: d.m}})
	}
}

func (cf *comFlow) addMsg(label string, value interface{}) {
	cf.sig.WriteLabel(label)
	cf.values = append(cf.values, value)
}

// EndValues ends the current state and commits its data to the parent state.
func (d *comFlow) EndValues() {
	sig := d.sig.String()
	if cnt := len(d.values); cnt == 0 {
		d.m.finishState(sig)
	} else {
		var v interface{}
		if cnt == 1 {
			v = d.values[0]
		} else {
			v = d.values
		}
		d.m.finishState(map[string]interface{}{
			sig: v,
		})
	}
}

// someone is trying to write a value to the flow, but we need keys to do that.
func (d *comFlow) commit(v interface{}) {
	d.m.Error(errutil.New("missing key when writing to a flow"))
}

// write the value into the key and change back to the flow state
func (d *comKey) commit(v interface{}) {
	d.comFlow.addMsg(d.key, v)
	d.m.changeState(d.comFlow)
}

// a literal exists within a flow, all flow functions except writing the literal value
// or ending the block result in an error.
// see also comBlock EndValues
func (d *comLiteral) EndValues() {
	d.m.finishState(d.out)
}

// a new value is being added to our slice
func (d *comSlice) commit(v interface{}) {
	d.values = append(d.values, v)
}

// the slice is done, write it to our parent whomever that is.
func (d *comSlice) EndValues() {
	d.m.finishState(d.values)
}

// compact raw values would normally just write the value
// but: we don't want to lose the *kind* of the choice
// so we do this specially
func (d *comSwap) SpecifyValue(kind string, value interface{}) {
	d.out = map[string]interface{}{
		kind + ":": value,
	}
}

func (d *comSwap) EndValues() {
	d.m.finishState(d.out)
}

// record the swap choice and move to an error detection state
func (d *comSwap) commit(v interface{}) {
	d.out = v
	d.m.changeState(&comWritten{d})
}

type comWritten struct {
	*comSwap
}

func (d *comWritten) commit(v interface{}) {
	d.m.Warning(errutil.New("value already committed"))
}
