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
	reg     []map[uint64]interface{}
	curMsg  json.RawMessage // can be set explicitly by the currently parsing context
	curFlow *cinFlow        // or, curMsg can be set via the current flow and key
	curKey  string
}

func Decode(dst jsn.Marshalee, msg json.RawMessage, reg []map[uint64]interface{}) error {
	dec := xDecoder{reg: reg, Machine: chart.MakeDecoder(custom)}
	next := dec.newBlock(msg, new(chart.StateMix))
	next.OnCommit = func(interface{}) {}
	dec.ChangeState(next)
	dst.Marshal(&dec)
	return dec.Errors()
}

func (dec *xDecoder) readMap(msg json.RawMessage) (okay bool) {
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

func (dec *xDecoder) readSlot(slot jsn.SlotBlock, msg json.RawMessage) (okay bool) {
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

func (dec *xDecoder) readFullSwap(p jsn.SwapBlock, msg json.RawMessage) (okay bool) {
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
			pick := newStringKey(sig.params[0].label)
			if ok := p.SetSwap(pick); !ok {
				dec.Error(errutil.Fmt("swap has unexpected choice %q", pick))
			} else {
				dec.PushState(dec.newSwap(args))
				okay = true
			}
		}
	}
	return
}

func (dec *xDecoder) readRepeat(slice jsn.SliceBlock, msg json.RawMessage) (okay bool) {
	var msgs []json.RawMessage
	if e := json.Unmarshal(msg, &msgs); e != nil {
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

func (dec *xDecoder) readValue(pv interface{}, msg json.RawMessage) (err error) {
	if el, ok := pv.(interface{ SetValue(interface{}) bool }); ok {
		var i interface{}
		if e := json.Unmarshal(msg, &i); e != nil {
			err = e
		} else if !el.SetValue(i) {
			err = errutil.New("couldnt set value", i)
		}
	} else {
		if e := json.Unmarshal(msg, pv); e != nil {
			err = e // couldn't unmarshal directly into the target value?
		}
	}
	dec.Commit("new value")
	return
}

func (dec *xDecoder) newBlock(msg json.RawMessage, next *chart.StateMix) *chart.StateMix {
	dec.curMsg = msg
	next.OnMap = func(_, _ string) bool {
		return dec.readMap(msg)
	}
	next.OnSlot = func(_ string, slot jsn.SlotBlock) bool {
		return dec.readSlot(slot, msg)
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	next.OnSwap = func(_ string, p jsn.SwapBlock) bool {
		return dec.readFullSwap(p, msg)
	}
	next.OnRepeat = func(_ string, slice jsn.SliceBlock) bool {
		return dec.readRepeat(slice, msg)
	}
	next.OnValue = func(_ string, pv interface{}) error {
		return dec.readValue(pv, msg)
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
	var next chart.StateMix
	next.OnKey = func(key, _ string) (err error) {
		dec.ChangeState(dec.newKeyValue(flow, key, next))
		return
	}
	next.OnEnd = func() {
		dec.FinishState("flow")
	}
	return &next
}

func (dec *xDecoder) GetCurrentMessage() (ret json.RawMessage) {
	if msg := dec.curMsg; msg != nil {
		ret = msg
	} else if flow, key := dec.curFlow, dec.curKey; flow != nil {
		if msg, e := flow.getArg(key); e != nil {
			dec.Error(e)
		} else {
			ret = msg
		}
		if ret != nil {
			dec.curMsg = ret // cache
		} else {
			dec.curMsg = emptyMessage // mark it
		}
	}
	return
}

var emptyMessage = make([]byte, 0, 0)

func (dec *xDecoder) ReadCurrentMessage(pv interface{}) (okay bool) {
	if msg := dec.GetCurrentMessage(); len(msg) > 0 {
		if e := json.Unmarshal(msg, pv); e == nil {
			okay = true
		}
	}
	return
}

// given a specific key, we have to look up the corresponding msg value before we can process it
// some states however treat that key differently... ie. embedded swaps
func (dec *xDecoder) newKeyValue(flow *cinFlow, key string, prev chart.StateMix) *chart.StateMix {
	next := prev // copy
	//
	dec.curFlow = flow
	dec.curKey = key
	dec.curMsg = nil
	//
	next.OnMap = func(_, _ string) (okay bool) {
		if msg := dec.GetCurrentMessage(); len(msg) > 0 {
			okay = dec.readMap(msg)
		}
		return
	}
	next.OnSlot = func(_ string, slot jsn.SlotBlock) (okay bool) {
		if msg := dec.GetCurrentMessage(); len(msg) > 0 {
			okay = dec.readSlot(slot, msg)
		}
		return
	}
	next.OnSwap = func(_ string, p jsn.SwapBlock) (okay bool) {
		if msg, pick, e := flow.getPick(key); e != nil {
			dec.Error(e)
		} else if len(msg) > 0 {
			if ok := p.SetSwap(newStringKey(pick)); !ok {
				dec.Error(errutil.Fmt("embedded swap at %q has unexpected choice %q", key, pick))
			} else {
				dec.PushState(dec.newSwap(msg))
				okay = true
			}
		}
		return
	}
	next.OnRepeat = func(_ string, slice jsn.SliceBlock) (okay bool) {
		if msg := dec.GetCurrentMessage(); len(msg) > 0 {
			okay = dec.readRepeat(slice, msg)
		}
		return
	}
	next.OnValue = func(_ string, pv interface{}) (err error) {
		if msg := dec.GetCurrentMessage(); len(msg) == 0 {
			err = jsn.Missing
		} else {
			err = dec.readValue(pv, msg)
		}
		return
	}
	next.OnCommit = func(interface{}) {
		dec.curFlow = nil
		dec.curKey = ""
		dec.curMsg = nil
		dec.ChangeState(&prev)
	}
	return &next
}

// we expect at least one msg, and no more and no fewer values than msgs
func (dec *xDecoder) newSlice(msgs []json.RawMessage) *chart.StateMix {
	next := dec.newBlock(msgs[0], new(chart.StateMix))
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
	var next chart.StateMix
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
	return &next
}

// pretty much any simple value or block data can fulfill a swap
func (dec *xDecoder) newSwap(msg json.RawMessage) *chart.StateMix {
	next := dec.newBlock(msg, new(chart.StateMix))
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
			} else if t, ok := findType(Hash(k), dec.reg); !ok {
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
		} else if t, ok := findType(Hash(k), dec.reg); !ok {
			err = errutil.New("couldnt find type for string", k)
		} else {
			retType, retKey, retVal = t, k, nil // no args
		}
	}
	return
}
