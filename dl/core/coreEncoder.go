package core

import (
	"encoding/json"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

func CompactEncoder(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
	typeName := flow.GetType()
	switch op := flow.GetFlow().(type) {
	case *CallPattern:
		err = EncodePattern(m, op)
	case *GetValue:
		if name := encodeAddress(op.Source); len(name) == 0 {
			err = chart.Unhandled(typeName)
		} else {
			err = m.MarshalValue(typeName, name)
		}
	case *VariableRef:
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
func encodeAddress(addr Address) (ret string) {
	if vref, ok := addr.Value.(*VariableRef); ok {
		ret = encodeVariableRef(vref)
	}
	return
}

func encodeVariableRef(vref *VariableRef) (ret string) {
	if len(vref.Dot) == 0 {
		if name, ok := vref.Name.(*literal.TextValue); ok {
			ret = "@" + name.Value
		}
	}
	return
}

func Decode(dst jsn.Marshalee, msg json.RawMessage, sig cin.Signatures) error {
	return cin.NewDecoder(sig).
		SetSlotDecoder(CompactSlotDecoder).
		Decode(dst, msg)
}

func CompactFlowDecoder(m jsn.Marshaler, flow jsn.FlowBlock, msg json.RawMessage) (err error) {
	var str string
	if e := json.Unmarshal(msg, &str); e != nil || len(str) == 0 && str[0] != '@' {
		err = chart.Unhandled(flow.GetType())
	} else {
		if (!flow.SetFlow(&VariableRef{Name: T(str[1:])})) {
			err = errutil.New("couldnt set variable ref")
		}
	}
	return
}

// unhandled reads are attempted via default readSlot evaluation.
func CompactSlotDecoder(m jsn.Marshaler, slot jsn.SlotBlock, msg json.RawMessage) (err error) {
	// switching on the slot ptr's type seems like it should work, but only results in untyped interfaces
	switch typeName := slot.GetType(); typeName {
	default:
		err = chart.Unhandled(typeName)
	case
		rt.BoolEval_Type,
		rt.NumberEval_Type,
		rt.TextEval_Type,
		rt.NumListEval_Type,
		rt.TextListEval_Type,
		rt.RecordEval_Type,
		rt.RecordListEval_Type:
		var str string
		if e := json.Unmarshal(msg, &str); e == nil && len(str) > 0 && str[0] == '@' {
			if !slot.SetSlot(GetVariable(str[1:])) {
				err = errutil.New("unexpected error setting slot")
			}
		} else {
			err = literal.CompactSlotDecoder(m, slot, msg)
		}
	}
	return
}
