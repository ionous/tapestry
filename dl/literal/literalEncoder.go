package literal

import (
	"encoding/json"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

func CompactEncoder(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
	typeName := flow.GetType()
	switch out := flow.GetFlow().(type) {
	default:
		err = chart.Unhandled(typeName)

	case *BoolValue:
		err = m.MarshalValue(typeName, out.Bool)

	case *NumValue:
		err = m.MarshalValue(typeName, out.Num)

	case *TextValue:
		err = m.MarshalValue(typeName, out.Text)

	case *NumValues:
		err = m.MarshalValue(typeName, out.Values)

	case *TextValues:
		err = m.MarshalValue(typeName, out.Values)
	}
	return
}

func CompactSlotDecoder(slot jsn.SlotBlock, msg json.RawMessage) (err error) {
	if ptr, e := readLiteral(slot.GetType(), "", msg); e != nil {
		err = e
	} else if !slot.SetSlot(ptr) {
		err = errutil.New("unexpected error setting slot")
	}
	return
}

func ReadLiteral(aff affine.Affinity, cls string, msg json.RawMessage) (ret LiteralValue, err error) {
	return readLiteral(aff.String()+"_eval", cls, msg)
}

func readLiteral(typeName, cls string, msg json.RawMessage) (ret LiteralValue, err error) {
	// switching on the slot ptr's type seems like it should work, but only results in untyped interfaces
	switch typeName {
	default:
		err = chart.Unhandled("CustomSlot")

	case rt.BoolEval_Type:
		var val bool
		if e := json.Unmarshal(msg, &val); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			ret = &BoolValue{Bool: val, Class: cls}
		}

	case rt.NumberEval_Type:
		var val float64
		if e := json.Unmarshal(msg, &val); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			ret = &NumValue{Num: val, Class: cls}
		}

	case rt.TextEval_Type:
		var val string
		if e := json.Unmarshal(msg, &val); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			ret = &TextValue{Text: val, Class: cls}
		}

	case rt.NumListEval_Type:
		var val []float64
		if e := json.Unmarshal(msg, &val); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			ret = &NumValues{Values: val, Class: cls}
		}

	case rt.TextListEval_Type:
		var val []string
		if e := json.Unmarshal(msg, &val); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			ret = &TextValues{Values: val, Class: cls}
		}
	}
	return
}
