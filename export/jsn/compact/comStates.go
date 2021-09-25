package compact

import (
	"strconv"

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
	name     string
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

var namer int

func newName(base string) string {
	return base + "-" + strconv.Itoa(namer)
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
		d.m.changeState(&comLiteral{comValue: comValue{m: d.m, name: newName("literal")}})
	}
}

func (cf *comFlow) addMsg(label string, value interface{}) {
	cf.sig.WriteLabel(label)
	cf.values = append(cf.values, value)
}

func (d *comFlow) readData() (ret interface{}) {
	sig := d.sig.String()
	if cnt := len(d.values); cnt == 0 {
		ret = sig
	} else {
		var v interface{}
		if cnt == 1 {
			v = d.values[0]
		} else {
			v = d.values
		}
		ret = map[string]interface{}{
			sig: v,
		}
	}
	return
}

func (d *comFlow) writeData(v interface{}) {
	d.m.Error(errutil.New("missing key when writing to a flow"))
}

// override the flow's name to prefix that we are a substate of it.
func (d *comKey) named() string {
	return "key-" + d.comFlow.named()
}

func (d *comKey) writeData(v interface{}) {
	d.comFlow.addMsg(d.key, v)
	d.m.changeState(d.comFlow)
}

// a literal exists within a flow, all flow functions except writing the literal value
// or ending the block result in an error.
// see also comBlock EndValues
func (d *comLiteral) EndValues() {
	was := d.m.popState()         // gets rid of us
	d.m.writeData(was.readData()) // write our accumulated data to the new parent
}

func (d *comSlice) readData() interface{} {
	return d.values
}

func (d *comSlice) writeData(v interface{}) {
	d.values = append(d.values, v)
}

// compact raw values would normally just write the value
// but: we don't want to lose the *kind* of the choice
// so we do this specially
func (d *comSwap) WriteValue(kind string, value interface{}) {
	d.m.writeData(map[string]interface{}{
		kind + ":": value,
	})
}

func (d *comSwap) writeData(v interface{}) {
	d.out = v
	d.m.changeState(&comWritten{d})
}

type comWritten struct {
	*comSwap
}

func (d *comWritten) writeData(v interface{}) {
	d.m.Warning(errutil.New(d.name, "already written"))
}
