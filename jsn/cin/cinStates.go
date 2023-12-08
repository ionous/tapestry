package cin

import (
	"errors"
	"unicode"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"github.com/ionous/errutil"
)

type xDecoder struct {
	chart.Machine
	TypeCreator
	customSlot SlotDecoder
}

// slots in memory are filled with commands; in scripts, however,
// they be shortcut values: like a string for a literal command.
// a slot decoder expands those shortcuts.
type SlotDecoder func(jsn.Marshaler, jsn.SlotBlock, any) error

// Customization of the decoding process
type Decoder xDecoder

func Decode(dst jsn.Marshalee, msg map[string]any, reg TypeCreator) error {
	dec := NewDecoder(reg)
	return dec.Decode(dst, msg)
}

func NewDecoder(reg TypeCreator) *Decoder {
	d := &xDecoder{
		TypeCreator: reg,
		Machine:     chart.MakeDecoder(),
		customSlot: func(jsn.Marshaler, jsn.SlotBlock, any) error {
			return chart.Unhandled("no custom slot handler")
		},
	}
	return (*Decoder)(d)
}

func (b *Decoder) SetSlotDecoder(handler SlotDecoder) *Decoder {
	x := (*xDecoder)(b)
	x.customSlot = handler
	return b
}

func (b *Decoder) Decode(dst jsn.Marshalee, msg map[string]any) error {
	x := (*xDecoder)(b)
	return x.Decode(dst, msg)
}

func (b *Decoder) DecodeValue(dst jsn.Marshalee, val any) error {
	x := (*xDecoder)(b)
	return x.Decode(dst, val)
}

func (dec *xDecoder) Decode(dst jsn.Marshalee, val any) error {
	next := dec.newBlock(val, new(chart.StateMix))
	next.OnCommit = func(interface{}) {}
	dec.ChangeState(next)
	dst.Marshal(dec)
	return dec.Errors()
}

func (dec *xDecoder) readFlow(flow jsn.FlowBlock, msg map[string]any) (okay bool) {
	if op, e := ReadOp(msg); e != nil {
		dec.Error(e)
	} else {
		if ptr := dec.Machine.Markout; ptr != nil {
			*ptr, dec.Machine.Markout = op.Markup, nil
		}
		if flowMsg, e := newFlowData(op); e != nil {
			dec.Error(e)
		} else if name := flow.GetLede(); name != flowMsg.name {
			dec.Error(errutil.Fmt("mismatched commands: expected %q have %q", name, flowMsg.name))
		} else {
			dec.PushState(dec.newFlow(flowMsg))
			okay = true
		}
	}
	return
}

func (dec *xDecoder) readSlot(slot jsn.SlotBlock, body any) (okay bool) {
	if e := dec.customSlot(dec, slot, body); e == nil {
		dec.Commit("customSlot")
	} else {
		var unhandled chart.Unhandled
		if !errors.As(e, &unhandled) {
			dec.customSlot(dec, slot, body)
			dec.Error(e)
		} else if msg, ok := body.(map[string]any); !ok {
			dec.Error(errutil.Fmt("expected a command message; got %T(%v)", body, body))
		} else if op, e := ReadOp(msg); e != nil {
			dec.Error(errutil.Fmt("couldnt read op %s %w", msg, e))
		} else if v, e := dec.NewFromSignature(Hash(op.Sig, slot.GetType())); e != nil {
			dec.Error(e)
		} else {
			//should we be doing this on "OnCommit" instead of before the contents of the slat have been read?
			//( right now the caller calls Marshal on the slat -- but we could do it in here... )
			if !slot.SetSlot(v) {
				dec.Error(errutil.Fmt("couldn't put %T into slot %T", v, slot))
			} else {
				dec.PushState(dec.newSlot(op))
				okay = true
			}
		}
	}
	return
}

func (dec *xDecoder) readFullSwap(p jsn.SwapBlock, msg map[string]any) (okay bool) {
	if op, e := ReadOp(msg); e != nil {
		dec.Error(e)
	} else if sig, args, e := op.ReadMsg(); e != nil {
		dec.Error(e)
	} else if len(sig.Params) != 1 || len(sig.Params[0].Choice) > 0 {
		dec.Error(errutil.New("expected exactly one choice in", op.Sig))
	} else {
		pick := newStringKey(sig.Params[0].Label)
		if ok := p.SetSwap(pick); !ok {
			dec.Error(errutil.Fmt("swap has unexpected choice %q", pick))
		} else {
			el := args[0]
			dec.PushState(dec.newSwap(el))
			okay = true
		}
	}
	return
}

// allows a single value, or an array of values
func (dec *xDecoder) readRepeat(tgt jsn.SliceBlock, val any) (okay bool) {
	var slice []any
	if vs, ok := val.([]any); ok {
		slice = vs
	} else {
		slice = []any{val}
	}
	// to distinguish b/t missing and empty repeats, set even if size zero.
	// note: we don't get here if the record was missing
	cnt := len(slice)
	var next *chart.StateMix
	tgt.SetSize(cnt)
	if cnt == 0 {
		next = chart.NewBlockResult(&dec.Machine, "empty slice")
	} else {
		next = dec.newSlice(slice)
	}
	dec.PushState(next)
	return true
}

func (dec *xDecoder) readValue(pout any, val any) (err error) {
	// hrm... "str" types dont use boxing so they dont use set value
	if el, ok := pout.(interface{ SetValue(any) bool }); ok {
		if !el.SetValue(val) {
			err = errutil.New("couldnt set value", val)
		}
	} else if pstr, isString := pout.(*string); !isString {
		err = errutil.Fmt("unexpected value of type %T", pout)
	} else {
		if str, ok := val.(string); ok {
			*pstr = str
		} else if slice, ok := val.([]any); ok {
			if lines, e := SliceLines(slice); e != nil {
				err = e
			} else {
				*pstr = lines // a string with newlines.
			}
		}
	}
	// even on error?
	dec.Commit("new value")
	return
}

// return a new state capable of decoding an incoming value
func (dec *xDecoder) newBlock(val any, next *chart.StateMix) *chart.StateMix {
	next.OnMap = func(_ string, flow jsn.FlowBlock) (okay bool) {
		if msg, ok := val.(map[string]any); ok {
			okay = dec.readFlow(flow, msg)
		}
		return
	}
	next.OnSlot = func(_ string, slot jsn.SlotBlock) bool {
		return dec.readSlot(slot, val)
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	next.OnSwap = func(_ string, swap jsn.SwapBlock) (okay bool) {
		if msg, ok := val.(map[string]any); ok {
			okay = dec.readFullSwap(swap, msg)
		}
		return
	}
	next.OnRepeat = func(_ string, slice jsn.SliceBlock) bool {
		return dec.readRepeat(slice, val)
	}
	next.OnValue = func(_ string, pv interface{}) error {
		return dec.readValue(pv, val)
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

// given a specific key, we have to look up the corresponding msg value before we can process it
// some states however treat that key differently... ie. embedded swaps
func (dec *xDecoder) newKeyValue(cmd *cinFlow, key string, prev chart.StateMix) *chart.StateMix {
	next := prev // copy
	//
	next.OnMap = func(_ string, flow jsn.FlowBlock) (okay bool) {
		if arg, e := cmd.getArg(key); e != nil {
			dec.Error(e)
		} else if val, ok := arg.(map[string]any); ok {
			okay = dec.readFlow(flow, val)
		}
		return
	}
	next.OnSlot = func(_ string, slot jsn.SlotBlock) (okay bool) {
		if arg, e := cmd.getArg(key); e != nil {
			dec.Error(e)
		} else if arg != nil {
			okay = dec.readSlot(slot, arg)
		}
		return
	}
	next.OnSwap = func(_ string, swap jsn.SwapBlock) (okay bool) {
		if arg, pick, e := cmd.getPick(key); e != nil {
			dec.Error(e)
		} else if arg != nil {
			if ok := swap.SetSwap(newStringKey(pick)); !ok {
				dec.Error(errutil.Fmt("embedded swap at %q has unexpected choice %q", key, pick))
			} else {
				dec.PushState(dec.newSwap(arg))
				okay = true
			}
		}
		return
	}
	next.OnRepeat = func(_ string, slice jsn.SliceBlock) (okay bool) {
		if arg, e := cmd.getArg(key); e != nil {
			dec.Error(e)
		} else if arg != nil {
			okay = dec.readRepeat(slice, arg)
		}
		return
	}
	next.OnValue = func(_ string, pv interface{}) (err error) {
		if arg, e := cmd.getArg(key); e != nil {
			dec.Error(e)
		} else if arg == nil {
			err = jsn.Missing
		} else {
			err = dec.readValue(pv, arg)
		}
		return
	}
	next.OnCommit = func(interface{}) {
		dec.ChangeState(&prev)
	}
	return &next
}

// assumes at least one value, and no more and no fewer values than msgs
func (dec *xDecoder) newSlice(slice []any) *chart.StateMix {
	next := dec.newBlock(slice[0], new(chart.StateMix))
	// loops by slicing the array each time.
	next.OnCommit = func(any) {
		var after *chart.StateMix
		if len(slice) > 1 {
			rest := slice[1:]
			after = dec.newSlice(rest)
		} else {
			after = chart.NewBlockResult(&dec.Machine, "slice end")
		}
		dec.ChangeState(after)
	}
	next.OnEnd = func() {
		dec.Error(errutil.Fmt("slice underflow, %d remaining", len(slice)))
		dec.FinishState("bad slice")
	}
	return next
}

// in compact format, the msg holds the slat which fills the slot
func (dec *xDecoder) newSlot(op Op) *chart.StateMix {
	var next chart.StateMix
	next.OnMap = func(string, jsn.FlowBlock) (okay bool) {
		if flow, e := newFlowData(op); e != nil {
			dec.Error(e)
		} else {
			if ptr := dec.Machine.Markout; ptr != nil {
				*ptr, dec.Machine.Markout = op.Markup, nil
			}
			dec.PushState(dec.newFlow(flow))
			okay = true
		}
		return
	}
	next.OnValue = func(typeName string, pv any) (err error) {
		if sig, args, e := op.ReadMsg(); e != nil {
			dec.Error(e)
		} else if cnt := len(args); cnt != 1 {
			dec.Error(errutil.Fmt("expected %s to have one arg, has %d", sig.Name, cnt))
		} else {
			err = dec.readValue(pv, args[0])
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
func (dec *xDecoder) newSwap(val any) *chart.StateMix {
	next := dec.newBlock(val, new(chart.StateMix))
	next.OnCommit = func(interface{}) {
		dec.ChangeState(chart.NewBlockResult(&dec.Machine, "swap end"))
	}
	next.OnEnd = func() {
		dec.FinishState("empty swap")
	}
	return next
}
