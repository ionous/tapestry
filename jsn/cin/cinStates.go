package cin

import (
	"encoding/json"
	"errors"
	"unicode"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"github.com/ionous/errutil"
)

type Signatures []map[uint64]interface{}

type xDecoder struct {
	chart.Machine
	reg        Signatures
	customSlot SlotDecoder
	customFlow FlowDecoder
}

type SlotDecoder func(jsn.Marshaler, jsn.SlotBlock, json.RawMessage) error
type FlowDecoder func(jsn.Marshaler, jsn.FlowBlock, json.RawMessage) error

// Customization of the decoding process
type Decoder xDecoder

func NewDecoder(reg Signatures) *Decoder {
	d := &xDecoder{
		reg:     reg,
		Machine: chart.MakeDecoder(),
		customSlot: func(jsn.Marshaler, jsn.SlotBlock, json.RawMessage) error {
			return chart.Unhandled("no custom slot handler")
		},
		customFlow: func(jsn.Marshaler, jsn.FlowBlock, json.RawMessage) error {
			return chart.Unhandled("no custom flow handler")
		},
	}
	return (*Decoder)(d)
}

func (b *Decoder) SetSlotDecoder(handler SlotDecoder) *Decoder {
	x := (*xDecoder)(b)
	x.customSlot = handler
	return b
}
func (b *Decoder) SetFlowDecoder(handler FlowDecoder) *Decoder {
	x := (*xDecoder)(b)
	x.customFlow = handler
	return b
}

func (b *Decoder) Decode(dst jsn.Marshalee, msg json.RawMessage) error {
	x := (*xDecoder)(b)
	return x.decode(dst, msg)
}

func Decode(dst jsn.Marshalee, msg json.RawMessage, reg Signatures) error {
	dec := NewDecoder(reg)
	return dec.Decode(dst, msg)
}

func (dec *xDecoder) decode(dst jsn.Marshalee, msg json.RawMessage) error {
	next := dec.newBlock(msg, new(chart.StateMix))
	next.OnCommit = func(interface{}) {}
	dec.ChangeState(next)
	dst.Marshal(dec)
	return dec.Errors()
}

func (dec *xDecoder) readFlow(flow jsn.FlowBlock, msg json.RawMessage) (okay bool) {
	if e := dec.customFlow(dec, flow, msg); e == nil {
		dec.Commit("customFlow")
	} else {
		var unhandled chart.Unhandled
		if !errors.As(e, &unhandled) {
			dec.Error(e)
		} else if cmd, e := dec.readCmd(msg); e != nil {
			dec.Error(e)
		} else {
			if dec.Machine.Comment != nil {
				*dec.Machine.Comment = cmd.comment
				dec.Machine.Comment = nil
			}
			if flow, e := newFlowData(cmd); e != nil {
				dec.Error(e)
			} else {
				dec.PushState(dec.newFlow(flow))
				okay = true
			}
		}
	}
	return
}

func (dec *xDecoder) readSlot(slot jsn.SlotBlock, msg json.RawMessage) (okay bool) {
	if e := dec.customSlot(dec, slot, msg); e == nil {
		dec.Commit("customSlot")
	} else {
		var unhandled chart.Unhandled
		if !errors.As(e, &unhandled) {
			dec.Error(e)
		} else if cmd, e := dec.readCmd(msg); e != nil {
			dec.Error(e)
		} else if v := newFromType(cmd.reg); !slot.SetSlot(v) {
			dec.Error(errutil.Fmt("couldn't put %T into slot %T", v, slot))
		} else {
			dec.PushState(dec.newSlot(cmd))
			okay = true
		}
	}
	return
}

func (dec *xDecoder) readFullSwap(p jsn.SwapBlock, msg json.RawMessage) (okay bool) {
	// expanded swaps { "swapName choice:": <value> }
	if cmd, e := dec.readCmd(msg); e != nil {
		dec.Error(e)
	} else if sig, e := cmd.getSignature(); e != nil {
		dec.Error(e)
	} else if len(sig.params) != 1 || len(sig.params[0].choice) > 0 {
		dec.Error(errutil.New("expected exactly one choice in", cmd.key))
	} else {
		pick := newStringKey(sig.params[0].label)
		if ok := p.SetSwap(pick); !ok {
			dec.Error(errutil.Fmt("swap has unexpected choice %q", pick))
		} else {
			dec.PushState(dec.newSwap(cmd.args))
			okay = true
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
	next.OnMap = func(_ string, flow jsn.FlowBlock) bool {
		return dec.readFlow(flow, msg)
	}
	next.OnSlot = func(_ string, slot jsn.SlotBlock) bool {
		return dec.readSlot(slot, msg)
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	next.OnSwap = func(_ string, swap jsn.SwapBlock) bool {
		return dec.readFullSwap(swap, msg)
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

var emptyMessage = make([]byte, 0, 0)

// given a specific key, we have to look up the corresponding msg value before we can process it
// some states however treat that key differently... ie. embedded swaps
func (dec *xDecoder) newKeyValue(msgs *cinFlow, key string, prev chart.StateMix) *chart.StateMix {
	next := prev // copy
	//
	next.OnMap = func(_ string, flow jsn.FlowBlock) (okay bool) {
		if msg, e := msgs.getArg(key); e != nil {
			dec.Error(e)
		} else if len(msg) > 0 {
			okay = dec.readFlow(flow, msg)
		}
		return
	}
	next.OnSlot = func(_ string, slot jsn.SlotBlock) (okay bool) {
		if msg, e := msgs.getArg(key); e != nil {
			dec.Error(e)
		} else if len(msg) > 0 {
			okay = dec.readSlot(slot, msg)
		}
		return
	}
	next.OnSwap = func(_ string, swap jsn.SwapBlock) (okay bool) {
		if msg, pick, e := msgs.getPick(key); e != nil {
			dec.Error(e)
		} else if len(msg) > 0 {
			if ok := swap.SetSwap(newStringKey(pick)); !ok {
				dec.Error(errutil.Fmt("embedded swap at %q has unexpected choice %q", key, pick))
			} else {
				dec.PushState(dec.newSwap(msg))
				okay = true
			}
		}
		return
	}
	next.OnRepeat = func(_ string, slice jsn.SliceBlock) (okay bool) {
		if msg, e := msgs.getArg(key); e != nil {
			dec.Error(e)
		} else if len(msg) > 0 {
			okay = dec.readRepeat(slice, msg)
		}
		return
	}
	next.OnValue = func(_ string, pv interface{}) (err error) {
		if msg, e := msgs.getArg(key); e != nil {
			dec.Error(e)
		} else if len(msg) == 0 {
			err = jsn.Missing
		} else {
			err = dec.readValue(pv, msg)
		}
		return
	}
	next.OnCommit = func(interface{}) {
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
func (dec *xDecoder) newSlot(cmd cmdData) *chart.StateMix {
	var next chart.StateMix
	next.OnMap = func(string, jsn.FlowBlock) (okay bool) {
		if flow, e := newFlowData(cmd); e != nil {
			dec.Error(e)
		} else {
			if dec.Machine.Comment != nil {
				*dec.Machine.Comment = cmd.comment
				dec.Machine.Comment = nil
			}
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

type cmdData struct {
	reg     interface{}
	key     string
	args    json.RawMessage
	comment string
}

func (c cmdData) getSignature() (ret sigReader, err error) {
	err = ret.readSig(c.key)
	return
}

// retType: type registry nil pointer
// retKey: raw signature in the json. ex. "Story:", "Always", or "Command some thing:else:"
// retVal: json msg data to the right of the key
func (dec *xDecoder) readCmd(msg json.RawMessage) (ret cmdData, err error) {
	// ex. {"Story:":  [...]}
	// except note that literals are stored as their value,
	// functions with one parameter dont use the array, and
	// functions without parameters are stored as simple strings.
	var d map[string]json.RawMessage
	if e := json.Unmarshal(msg, &d); e == nil {
		var out cmdData
		for k, args := range d {
			if k == commentMarker {
				if e := json.Unmarshal(args, &out.comment); e != nil {
					err = errutil.New("couldnt read comment", e)
					break
				}
			} else if len(out.key) > 0 {
				err = errutil.New("expected only a single key", d)
				break
			} else {
				out.key, out.args = k, args
				continue // keep going to catch errors
			}
		}
		if err == nil {
			if len(out.comment) > 0 && len(out.key) == 0 {
				out.key = commentMarker
			}
			if t, ok := findType(Hash(out.key), dec.reg); !ok {
				err = errutil.New("couldnt find type for signature", out.key)
			} else {
				out.reg = t
				ret = out
			}
		}
	} else {
		var k string // parameterless commands can be simple strings.
		if e := json.Unmarshal(msg, &k); e != nil {
			err = e
		} else if t, ok := findType(Hash(k), dec.reg); !ok {
			err = errutil.New("couldnt find type for string", k)
		} else {
			ret.reg, ret.key = t, k // no args, its parameterless
		}
	}
	return
}

const commentMarker = "--"
