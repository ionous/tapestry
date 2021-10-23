package din

import (
	"encoding/json"

	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"github.com/ionous/errutil"
)

type xDecoder struct {
	chart.Machine
	reg Registry
}

func Decode(dst jsn.Marshalee, reg Registry, msg json.RawMessage) error {
	dec := xDecoder{reg: reg}
	next := dec.newBlock(&msg)
	next.OnCommit = func(interface{}) {}
	dec.ChangeState(next)
	dst.Marshal(&dec)
	return dec.Errors()
}

func (dec *xDecoder) newValue(pm *json.RawMessage, next *chart.StateMix) *chart.StateMix {
	next.OnValue = func(typeName string, pv interface{}) {
		var d dinValue
		if e := json.Unmarshal(*pm, &d); e != nil {
			dec.Error(e)
		} else {
			if el, ok := pv.(interface{ SetValue(interface{}) bool }); ok {
				var i interface{}
				if e := json.Unmarshal(d.Msg, &i); e != nil {
					dec.Error(e)
				} else if !el.SetValue(i) {
					dec.Error(errutil.New("couldnt set value", i))
				}
			} else {
				if e := json.Unmarshal(d.Msg, pv); e != nil {
					dec.Error(e) // couldnt unmarshal directly into the target value
				}
			}
		}
		dec.Commit("new value")
	}
	// next.OnCommit -- handled by each caller
	return next
}

func (dec *xDecoder) newBlock(pm *json.RawMessage) *chart.StateMix {
	return dec.addBlock(pm, chart.NewReportingState(&dec.Machine))
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
	next.OnSlot = func(typeName string, slot jsn.Spotter) (okay bool) {
		var d, inner dinValue
		if e := json.Unmarshal(*pm, &d); e != nil {
			dec.Error(e)
		} else if d.Type != typeName {
			dec.Error(errutil.New("expected", typeName, "found", d.Type))
		} else if e := json.Unmarshal(d.Msg, &inner); e != nil {
			dec.Error(e)
		} else if v, e := dec.reg.NewType(inner.Type); e != nil {
			dec.Error(e)
		} else if !slot.SetSlot(v) {
			dec.Error(errutil.Fmt("couldn't put %T into slot %T", v, slot))
		} else {
			dec.PushState(dec.newSlot(d.Msg))
			okay = true
		}
		return
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	next.OnPick = func(typeName string, p jsn.Picker) (okay bool) {
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
				} else if _, ok := p.SetChoice(k); !ok {
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
	next.OnRepeat = func(_ string, slice jsn.Slicer) (okay bool) {
		var vs []json.RawMessage
		if e := json.Unmarshal(*pm, &vs); e != nil {
			dec.Error(e)
		} else if cnt := len(vs); cnt > 0 {
			slice.SetSize(cnt)
			dec.PushState(dec.newSlice(vs))
			okay = true
		}
		return
	}
	return next
}

func (dec *xDecoder) newFlow(fields dinFields) *chart.StateMix {
	next := chart.NewReportingState(&dec.Machine)
	next.OnKey = func(_, key string) (okay bool) {
		if msg, ok := fields[key]; ok {
			dec.ChangeState(dec.newKey(*next, msg))
			okay = true
		}
		return okay
	}
	next.OnEnd = func() {
		dec.FinishState(nil)
	}
	return next
}

func (dec *xDecoder) newKey(prev chart.StateMix, msg json.RawMessage) *chart.StateMix {
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
