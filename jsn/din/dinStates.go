package din

import (
	"encoding/json"
	r "reflect"

	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"github.com/ionous/errutil"
)

type xDecoder struct {
	chart.Machine
	reg map[string]r.Type
}

func Decode(dst jsn.Marshalee, reg map[string]r.Type, msg json.RawMessage) error {
	dec := xDecoder{reg: reg}
	next := dec.newBlock(&msg)
	next.OnCommit = func(interface{}) {}
	dec.ChangeState(next)
	dst.Marshal(&dec)
	return dec.Errors()
}

func (dec *xDecoder) newType(typeName string) (ret interface{}, err error) {
	if rtype, ok := dec.reg[typeName]; !ok {
		err = errutil.New("unknown type", typeName)
	} else {
		ret = r.New(rtype).Interface()
	}
	return
}

func (dec *xDecoder) newValue(pm *json.RawMessage, next *chart.StateMix) *chart.StateMix {
	next.OnValue = func(typeName string, pv interface{}) (err error) {
		var d dinValue
		if e := json.Unmarshal(*pm, &d); e != nil {
			err = e
		} else {
			if el, ok := pv.(interface{ SetValue(interface{}) bool }); ok {
				var i interface{}
				if e := json.Unmarshal(d.Msg, &i); e != nil {
					err = e
				} else if !el.SetValue(i) {
					err = errutil.New("couldnt set value", i)
				}
			} else {
				if e := json.Unmarshal(d.Msg, pv); e != nil {
					err = e // couldnt unmarshal directly into the target value
				}
			}
		}
		dec.Commit("new value")
		return
	}
	// next.OnCommit -- handled by each caller
	return next
}

func (dec *xDecoder) newBlock(pm *json.RawMessage) *chart.StateMix {
	var next chart.StateMix
	return dec.addBlock(pm, &next)
}

func (dec *xDecoder) addBlock(pm *json.RawMessage, next *chart.StateMix) *chart.StateMix {
	next.OnMap = func(_, typeName string) (okay bool) {
		var d dinMap
		if e := json.Unmarshal(*pm, &d); e != nil {
			dec.Error(e)
		} else if d.Type != typeName {
			dec.Error(errutil.New("expected", typeName, "found", d.Type))
		} else {
			dec.PushState(dec.newFlow(d.Fields))
			okay = true
		}
		return
	}
	next.OnSlot = func(typeName string, slot jsn.SlotBlock) (okay bool) {
		var d, inner dinValue
		if e := json.Unmarshal(*pm, &d); e != nil {
			dec.Error(e)
		} else if d.Type != typeName {
			dec.Error(errutil.New("expected", typeName, "found", d.Type))
		} else if e := json.Unmarshal(d.Msg, &inner); e != nil {
			dec.Error(e)
		} else if v, e := dec.newType(inner.Type); e != nil {
			dec.Error(e)
		} else if !slot.SetSlot(v) {
			dec.Error(errutil.Fmt("couldn't put %T into slot %T", v, slot))
		} else if v != nil {
			dec.PushState(dec.newSlot(d.Msg))
			okay = true
		}
		return
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	next.OnSwap = func(typeName string, p jsn.SwapBlock) (okay bool) {
		var d dinMap
		if e := json.Unmarshal(*pm, &d); e != nil {
			dec.Error(e)
		} else if d.Type != typeName {
			dec.Error(errutil.New("expected", typeName, "found", d.Type))
		} else {
			for k, v := range d.Fields {
				if okay {
					dec.Error(errutil.New("swap has too many choices"))
					break
				} else if ok := p.SetSwap(k); !ok {
					dec.Error(errutil.New("swap has unexpected choice", k))
					break
				} else {
					dec.PushState(dec.newSwap(v))
					okay = true
				}
			}
		}
		return
	}
	next.OnRepeat = func(_ string, slice jsn.SliceBlock) (okay bool) {
		var msgs []json.RawMessage
		if e := json.Unmarshal(*pm, &msgs); e != nil {
			dec.Error(e)
		} else {
			// to distinguish b/t missing and empty repeats, set even if size zero.
			// note: we don't get here if the record was missing
			var next *chart.StateMix
			cnt := len(msgs)
			slice.SetSize(cnt)
			if cnt == 0 {
				next = chart.NewBlockResult(&dec.Machine, "empty slice")
			} else {
				next = dec.newSlice(msgs)
			}
			dec.PushState(next)
			okay = true
		}
		return
	}
	return next
}

func (dec *xDecoder) newFlow(fields dinFields) *chart.StateMix {
	var next chart.StateMix
	next.OnKey = func(_, key string) (err error) {
		if msg, ok := fields[key]; !ok {
			err = jsn.Missing
		} else {
			dec.ChangeState(dec.newKeyValue(next, msg))
		}
		return
	}
	next.OnEnd = func() {
		dec.FinishState(nil)
	}
	return &next
}

func (dec *xDecoder) newKeyValue(prev chart.StateMix, msg json.RawMessage) *chart.StateMix {
	// a key's value can be a simple value, or a block.
	next := dec.newValue(&msg, dec.addBlock(&msg, &prev))
	next.OnCommit = func(interface{}) {
		dec.ChangeState(&prev)
	}
	return next
}

// we expect to get no more and no fewer values than msgs
func (dec *xDecoder) newSlice(msgs []json.RawMessage) *chart.StateMix {
	msg := msgs[0]
	next := dec.newValue(&msg, dec.newBlock(&msg))
	next.OnCommit = func(interface{}) {
		if msgs = msgs[1:]; len(msgs) > 0 {
			msg = msgs[0]
		} else {
			// expect nothing else now...
			dec.ChangeState(chart.NewBlockResult(&dec.Machine, nil))
		}
	}
	next.OnEnd = func() {
		dec.Error(errutil.New("slice underflow", len(msgs)))
		dec.FinishState(nil)
	}
	return next
}

func (dec *xDecoder) newSlot(msg json.RawMessage) *chart.StateMix {
	next := dec.newValue(&msg, dec.newBlock(&msg))
	next.OnCommit = func(interface{}) {
		dec.ChangeState(chart.NewBlockResult(&dec.Machine, nil))
	}
	next.OnEnd = func() {
		dec.FinishState(nil)
	}
	return next
}

func (dec *xDecoder) newSwap(msg json.RawMessage) *chart.StateMix {
	next := dec.newValue(&msg, dec.newBlock(&msg))
	next.OnCommit = func(interface{}) {
		dec.ChangeState(chart.NewBlockResult(&dec.Machine, nil))
	}
	next.OnEnd = func() {
		dec.FinishState(nil)
	}
	return next
}
