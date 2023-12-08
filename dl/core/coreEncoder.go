package core

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
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

// unhandled reads are attempted via default readSlot evaluation.
func CompactSlotDecoder(m jsn.Marshaler, slot jsn.SlotBlock, body any) (err error) {
	// switching on the slot ptr's type seems like it should work, but only results in untyped interfaces
	switch typeName := slot.GetType(); typeName {
	default:
		err = chart.Unhandled(typeName)
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
			if !slot.SetSlot(Variable(str)) {
				err = errutil.New("unexpected error setting slot")
			}
		} else {
			err = literal.CompactSlotDecoder(m, slot, body)
		}
	}
	return
}

func getVariableString(val any) (ret string) {
	if str, ok := val.(string); ok && len(str) > 0 && str[0] == '@' {
		ret = str[1:]
	}
	return
}
