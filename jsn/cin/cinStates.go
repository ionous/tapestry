package cin

import (
	"errors"
	r "reflect"
	"unicode"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"github.com/ionous/errutil"
)

type xDecoder struct {
	chart.Machine
	TypeCreator
	customSlot SlotDecoder
	customFlow FlowDecoder
}

type SlotDecoder func(jsn.Marshaler, jsn.SlotBlock, r.Value) error
type FlowDecoder func(jsn.Marshaler, jsn.FlowBlock, r.Value) error

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
		customSlot: func(jsn.Marshaler, jsn.SlotBlock, r.Value) error {
			return chart.Unhandled("no custom slot handler")
		},
		customFlow: func(jsn.Marshaler, jsn.FlowBlock, r.Value) error {
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

func (b *Decoder) Decode(dst jsn.Marshalee, msg map[string]any) error {
	x := (*xDecoder)(b)
	return x.Decode(dst, r.ValueOf(msg))
}

func (b *Decoder) DecodeValue(dst jsn.Marshalee, msg r.Value) error {
	x := (*xDecoder)(b)
	return x.Decode(dst, msg)
}

func (dec *xDecoder) Decode(dst jsn.Marshalee, msg r.Value) error {
	next := dec.newBlock(msg, new(chart.StateMix))
	next.OnCommit = func(interface{}) {}
	dec.ChangeState(next)
	dst.Marshal(dec)
	return dec.Errors()
}

func (dec *xDecoder) readFlow(flow jsn.FlowBlock, msg r.Value) (okay bool) {
	if e := dec.customFlow(dec, flow, msg); e == nil {
		dec.Commit("customFlow")
	} else {
		var unhandled chart.Unhandled
		if !errors.As(e, &unhandled) {
			dec.Error(e)
		} else if op, e := ReadOp(msg); e != nil {
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
	}
	return
}

func (dec *xDecoder) readSlot(slot jsn.SlotBlock, msg r.Value) (okay bool) {
	if e := dec.customSlot(dec, slot, msg); e == nil {
		dec.Commit("customSlot")
	} else {
		var unhandled chart.Unhandled
		if !errors.As(e, &unhandled) {
			dec.customSlot(dec, slot, msg)
			dec.Error(e)
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

func (dec *xDecoder) readFullSwap(p jsn.SwapBlock, msg r.Value) (okay bool) {
	// expanded swaps { "swapName choice:": <value> }
	if t := msg.Type(); !IsValidMap(t) {
		e := errutil.Fmt("expected a slot, have %s", t)
		dec.Error(e)
	} else if op, e := ReadOp(msg); e != nil {
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
			el := args.Index(0).Elem()
			dec.PushState(dec.newSwap(el))
			okay = true
		}
	}
	return
}

// allows a single value, or an array of values
func (dec *xDecoder) readRepeat(tgt jsn.SliceBlock, msg r.Value) (okay bool) {
	var slice r.Value
	if IsValidSlice(msg.Type()) {
		slice = msg
	} else {
		slice = r.ValueOf([]any{msg.Interface()})
	}
	// to distinguish b/t missing and empty repeats, set even if size zero.
	// note: we don't get here if the record was missing
	cnt := slice.Len()
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

func (dec *xDecoder) readValue(pv any, msg r.Value) (err error) {
	// hrm... "str" types dont use boxing so they dont use set value
	if el, ok := pv.(interface{ SetValue(any) bool }); ok {
		if !el.SetValue(msg.Interface()) {
			err = errutil.New("couldnt set value", msg)
		}
	} else if pstr, isString := pv.(*string); !isString {
		err = errutil.Fmt("unexpected value of type %T", pv)
	} else {
		if msg.Kind() == r.String {
			*pstr = msg.String()
		} else if lines, e := SliceLines(msg); e != nil {
			err = e
		} else {
			*pstr = lines // a string with newlines.
		}
	}
	// even on error?
	dec.Commit("new value")
	return
}

// allows any r.Value
func (dec *xDecoder) newBlock(msg r.Value, next *chart.StateMix) *chart.StateMix {
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

// given a specific key, we have to look up the corresponding msg value before we can process it
// some states however treat that key differently... ie. embedded swaps
func (dec *xDecoder) newKeyValue(cmd *cinFlow, key string, prev chart.StateMix) *chart.StateMix {
	next := prev // copy
	//
	next.OnMap = func(_ string, flow jsn.FlowBlock) (okay bool) {
		if arg, e := cmd.getArg(key); e != nil {
			dec.Error(e)
		} else if arg.IsValid() {
			okay = dec.readFlow(flow, arg)
		}
		return
	}
	next.OnSlot = func(_ string, slot jsn.SlotBlock) (okay bool) {
		if arg, e := cmd.getArg(key); e != nil {
			dec.Error(e)
		} else if arg.IsValid() {
			okay = dec.readSlot(slot, arg)
		}
		return
	}
	next.OnSwap = func(_ string, swap jsn.SwapBlock) (okay bool) {
		if arg, pick, e := cmd.getPick(key); e != nil {
			dec.Error(e)
		} else if arg.IsValid() {
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
		} else if arg.IsValid() {
			okay = dec.readRepeat(slice, arg)
		}
		return
	}
	next.OnValue = func(_ string, pv interface{}) (err error) {
		if arg, e := cmd.getArg(key); e != nil {
			dec.Error(e)
		} else if !arg.IsValid() {
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

// we expect at least one msg, and no more and no fewer values than msgs
func (dec *xDecoder) newSlice(slice r.Value) *chart.StateMix {
	first := slice.Index(0).Elem()
	next := dec.newBlock(first, new(chart.StateMix))
	// loops by slicing the array each time.
	next.OnCommit = func(interface{}) {
		var after *chart.StateMix
		if cnt := slice.Len(); cnt > 1 {
			rest := slice.Slice(1, cnt)
			after = dec.newSlice(rest)
		} else {
			after = chart.NewBlockResult(&dec.Machine, "slice end")
		}
		dec.ChangeState(after)
	}
	next.OnEnd = func() {
		dec.Error(errutil.Fmt("slice underflow, %d remaining", slice.Len()))
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
		} else if cnt := args.Len(); cnt != 1 {
			dec.Error(errutil.Fmt("expected %s to have one arg, has %d", sig.Name, cnt))
		} else {
			el := args.Index(0).Elem()
			err = dec.readValue(pv, el)
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
func (dec *xDecoder) newSwap(msg r.Value) *chart.StateMix {
	next := dec.newBlock(msg, new(chart.StateMix))
	next.OnCommit = func(interface{}) {
		dec.ChangeState(chart.NewBlockResult(&dec.Machine, "swap end"))
	}
	next.OnEnd = func() {
		dec.FinishState("empty swap")
	}
	return next
}
