package core

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// move some part of this into package assign
func CompactEncoder(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
	typeName := flow.GetType()
	switch op := flow.GetFlow().(type) {
	case *assign.CallPattern:
		err = assign.EncodePattern(m, op)

	case *assign.VariableRef:
		if name := encodeVariableRef(op); len(name) == 0 {
			err = chart.Unhandled(typeName)
		} else {
			err = m.MarshalValue(typeName, name)
		}

	case *literal.TextValue:
		// if the text starts with an @, skip it:
		// ( ie. dont confuse the rare text literal starting with an ampersand, with GetVar )
		if str := op.Value; len(str) > 0 && str[0] == '@' {
			err = chart.Unhandled(typeName)
		} else {
			err = literal.CompactEncoder(m, flow)
		}

	default:
		err = literal.CompactEncoder(m, flow)
	}
	return
}

func CustomEncoder(enc *encode.Encoder, op any) (ret any, err error) {
	switch op := op.(type) {
	case *assign.CallPattern:
		ret, err = assign.CustomEncoder(enc, op)

	case *assign.VariableRef:
		if name := encodeVariableRef(op); len(name) == 0 {
			err = compact.Unhandled
		} else {
			ret = name
		}

	case *literal.TextValue:
		// if the text starts with an @, skip it:
		// ( ie. dont confuse the rare text literal starting with an ampersand, with GetVar )
		if str := op.Value; len(str) > 0 && str[0] == '@' {
			err = compact.Unhandled
		} else {
			ret, err = literal.CustomEncoder(enc, op)
		}

	default:
		ret, err = literal.CustomEncoder(enc, op)
	}
	return
}

// write variables as a string prepended by @
// fix: it'd be nice if all parts were literals to write dot/bracket syntax a.b[5]
// fix: it'd be nicest if this could use package express to handle the parsing.
func encodeVariableRef(vref *assign.VariableRef) (ret string) {
	if len(vref.Dot) == 0 {
		if name, ok := vref.Name.(*literal.TextValue); ok {
			ret = "@" + name.Value
		}
	}
	return
}

func Decode(dst jsn.Marshalee, msg map[string]any, sig cin.Signatures) error {
	return DecodeValue(dst, msg, sig)
}

func DecodeValue(dst jsn.Marshalee, val any, sig cin.Signatures) error {
	return cin.NewDecoder(sig).
		SetSlotDecoder(CompactSlotDecoder).
		DecodeValue(dst, val)
}

func CoreDecoder(dec *decode.Decoder, slot string, body any) (ret any, err error) {
	// switching on the slot ptr's type seems like it should work, but only results in untyped interfaces
	switch slot {
	default:
		err = compact.Unhandled
	case
		// reading from a variable:
		rt.BoolEval_Type,
		rt.NumberEval_Type,
		rt.TextEval_Type,
		rt.NumListEval_Type,
		rt.TextListEval_Type,
		rt.RecordEval_Type,
		rt.RecordListEval_Type,
		// writing to a variable:
		assign.Address_Type:
		if str := getVariableString(body); len(str) > 0 {
			ret = Variable(str)
		} else {
			ret, err = literal.LiteralDecoder(dec, slot, body)
		}
	}
	return
}

// unhandled reads are attempted via default readSlot evaluation.
func CompactSlotDecoder(m jsn.Marshaler, slot jsn.SlotBlock, body any) (err error) {
	if v, e := CoreDecoder(nil, slot.GetType(), body); e != nil {
		if e == compact.Unhandled {
			err = chart.Unhandled("compact")
		} else {
			err = e
		}
	} else if !slot.SetSlot(v) {
		err = errutil.New("unexpected error setting slot")
	}
	return
}

func getVariableString(val any) (ret string) {
	if str, ok := val.(string); ok && len(str) > 0 && str[0] == '@' {
		ret = str[1:]
	}
	return
}
