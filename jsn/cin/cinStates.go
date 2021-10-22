package cin

import (
	"encoding/json"
	"unicode"

	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"github.com/ionous/errutil"
)

type xDecoder struct {
	chart.Machine
	reg            []map[uint64]interface{}
	CurrentMessage json.RawMessage
}

func Decode(dst jsn.Marshalee, msg json.RawMessage, reg []map[uint64]interface{}) error {
	dec := xDecoder{reg: reg, Machine: chart.MakeDecoder(custom)}
	next := dec.newBlock(msg)
	next.OnCommit = func(interface{}) {}
	dec.ChangeState(next)
	dst.Marshal(&dec)
	return dec.Errors()
}

func (dec *xDecoder) newValue(msg json.RawMessage, next *chart.StateMix) *chart.StateMix {
	dec.CurrentMessage = msg
	next.OnValue = func(_ string, pv interface{}) {
		if el, ok := pv.(interface{ SetValue(interface{}) bool }); ok {
			var i interface{}
			if e := json.Unmarshal(msg, &i); e != nil {
				dec.Error(e)
			} else if !el.SetValue(i) {
				dec.Error(errutil.New("couldnt set value", i))
			}
		} else {
			if e := json.Unmarshal(msg, pv); e != nil {
				dec.Error(e) // coudn't unmarshal directly into the target value?
			}
		}
		dec.Commit("new value")
	}
	return next
}

func (dec *xDecoder) newBlock(msg json.RawMessage) *chart.StateMix {
	return dec.addBlock(msg, chart.NewReportingState(&dec.Machine))
}

func (dec *xDecoder) addBlock(msg json.RawMessage, next *chart.StateMix) *chart.StateMix {
	next.OnMap = func(lede, _ string) (okay bool) {
		if _, k, args, e := dec.readCmd(msg); e != nil {
			dec.Error(e)
		} else if flow, e := newFlowData(k, args); e != nil {
			dec.Error(e)
		} else {
			dec.PushState(dec.newFlow(flow))
			okay = true
		}
		return
	}
	next.OnSlot = func(_ string, slot jsn.Spotter) (okay bool) {
		if t, k, args, e := dec.readCmd(msg); e != nil {
			dec.Error(e)
		} else if v := newFromType(t); !slot.SetSlot(v) {
			dec.Error(errutil.Fmt("couldn't put %T into slot %T", v, slot))
		} else {
			dec.PushState(dec.newSlot(k, args))
			okay = true
		}
		return
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	next.OnPick = func(_ string, p jsn.Picker) (okay bool) {
		// expanded swaps { "swapName choice:": <value> }
		if _, k, args, e := dec.readCmd(msg); e != nil {
			dec.Error(e)
		} else {
			// fix? if we cast the type to composer we could compare to typeName
			var sig sigReader
			if e := sig.readSig(k); e != nil {
				dec.Error(e)
			} else if len(sig.params) != 1 || len(sig.params[0].choice) > 0 {
				dec.Error(errutil.New("expected exactly one choice in", k))
			} else {
				choice := sig.params[0].label
				if _, ok := p.SetChoice(choice); !ok {
					dec.Error(errutil.New("swap has unexpected choice", k))
				} else {
					dec.PushState(dec.newSwap(args))
					okay = true
				}
			}
		}
		return
	}
	next.OnRepeat = func(_ string, slice jsn.Slicer) (okay bool) {
		var vs []json.RawMessage
		if e := json.Unmarshal(msg, &vs); e != nil {
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

// the message data is special, and the next state is expected to be a swap
func (dec *xDecoder) newEmbeddedSwap(prev chart.StateMix, msg json.RawMessage, pick string) *chart.StateMix {
	next := chart.NewReportingState(&dec.Machine)
	next.OnPick = func(typeName string, p jsn.Picker) (okay bool) {
		pick := newStringKey(pick)
		if _, ok := p.SetChoice(pick); !ok {
			dec.Error(errutil.Fmt("swap has unexpected %q", pick))
		} else {
			dec.PushState(dec.newSwap(msg))
			okay = true
		}
		return
	}
	next.OnCommit = func(interface{}) {
		dec.ChangeState(&prev)
	}
	return next
}

func newStringKey(s string) string {
	rs := make([]rune, 0, len(s)+1)
	rs = append(rs, '$')
	for _, r := range s {
		rs = append(rs, unicode.ToUpper(r))
	}
	return string(rs)
}

func (dec *xDecoder) newFlow(flow *cinFlow) *chart.StateMix {
	next := chart.NewReportingState(&dec.Machine)
	// the generated code is going to be calling this zero or more times
	// we need to walk the parameter names in order looking for matches
	next.OnKey = func(lede, _ string) (okay bool) {
		pick, msg := flow.findArg(lede)
		if len(pick) > 0 {
			dec.ChangeState(dec.newEmbeddedSwap(*next, msg, pick))
			okay = true
		} else if len(msg) > 0 {
			dec.ChangeState(dec.newKey(*next, msg))
			okay = true
		}
		return okay
	}
	next.OnEnd = func() {
		dec.FinishState("flow")
	}
	return next
}

// we have a valid flow key, we are now waiting on its value.
func (dec *xDecoder) newKey(prev chart.StateMix, msg json.RawMessage) *chart.StateMix {
	next := dec.newValue(msg, dec.addBlock(msg, &prev))
	next.OnCommit = func(interface{}) {
		dec.ChangeState(&prev)
	}
	return next
}

// we expect at least one msg, and no more and no fewer values than msgs
func (dec *xDecoder) newSlice(msgs []json.RawMessage) *chart.StateMix {
	next := dec.newValue(msgs[0], dec.newBlock(msgs[0]))
	next.OnCommit = func(interface{}) {
		var after *chart.StateMix
		if rest := msgs[1:]; len(rest) > 0 {
			after = dec.newSlice(rest)
		} else {
			after = chart.NewBlockResult(&dec.Machine, "slice end")
		}
		dec.ChangeState(after)
	}
	next.OnEnd = func() {
		dec.Error(errutil.New("slice underflow", len(msgs)))
		dec.FinishState("bad slice")
	}
	return next
}

// in compact format, the msg holds the slat which fills the slot
func (dec *xDecoder) newSlot(k string, args json.RawMessage) *chart.StateMix {
	next := chart.NewReportingState(&dec.Machine)
	next.OnMap = func(_, _ string) (okay bool) {
		if flow, e := newFlowData(k, args); e != nil {
			dec.Error(e)
		} else {
			dec.PushState(dec.newFlow(flow))
			okay = true
		}
		return
	}
	next.OnCommit = func(interface{}) {
		dec.ChangeState(chart.NewBlockResult(&dec.Machine, "slot end"))
	}
	next.OnEnd = func() {
		dec.FinishState("empty slot")
	}
	return next
}

// pretty much any simple value or block data can fulfill a swap
func (dec *xDecoder) newSwap(msg json.RawMessage) *chart.StateMix {
	next := dec.newValue(msg, dec.newBlock(msg))
	next.OnCommit = func(interface{}) {
		dec.ChangeState(chart.NewBlockResult(&dec.Machine, "swap end"))
	}
	next.OnEnd = func() {
		dec.FinishState("empty swap")
	}
	return next
}

// retType: type registry nil pointer
// retKey: raw signature in the json. ex. "Story:", "Always", or "Command some thing:else:"
// retVal: json msg data to the right of the key
func (dec *xDecoder) readCmd(msg json.RawMessage) (retType interface{}, retKey string, retVal json.RawMessage, err error) {
	// ex. {"Story:":  [...]}
	// except note that literals are stored as their value,
	// functions with one parameter dont use the array, and
	// functions without parameters are stored as simple strings.
	var d map[string]json.RawMessage
	if e := json.Unmarshal(msg, &d); e == nil {
		var found bool
		for k, args := range d {
			if found {
				err = errutil.New("expected only a single key", d)
				break
			} else if t, ok := findType(hash(k), dec.reg); !ok {
				err = errutil.New("couldnt find type for field", k)
				break
			} else {
				retType, retKey, retVal = t, k, args
				found = true
				continue // keep going to catch errors
			}
		}
		if err == nil && !found {
			err = errutil.New("expected a valid signature", d)
		}
	} else {
		var k string // parameterless commands can be simple strings.
		if e := json.Unmarshal(msg, &k); e != nil {
			err = e
		} else if t, ok := findType(hash(k), dec.reg); !ok {
			err = errutil.New("couldnt find type for string", k)
		} else {
			retType, retKey, retVal = t, k, nil // no args
		}
	}
	return
}
