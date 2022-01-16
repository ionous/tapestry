package literal

import (
	"encoding/json"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/rt"
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

	// records dont want to contain the names of their records
	// since that's generally recoverable from knowing the names of the fields
	// it might be better to allow that though, and only remove that info on the special compact encoding
	// otherwise writing are reading arent exactly idempotent ( writes from fields but generates a record )
	case *FieldValues:
		obj := make(map[string]interface{})
		if e := marshalFields(obj, out.Contains); e != nil {
			err = e
		} else {
			err = m.MarshalValue(typeName, obj)
		}
	}
	return
}

func CompactSlotDecoder(m jsn.Marshaler, slot jsn.SlotBlock, msg json.RawMessage) (err error) {
	if ptr, e := readLiteral(slot.GetType(), "", msg); e != nil {
		err = e
	} else if !slot.SetSlot(ptr) {
		err = errutil.New("unexpected error setting slot")
	}
	return
}

func ReadLiteral(aff affine.Affinity, cls string, msg json.RawMessage) (ret LiteralValue, err error) {
	// most literals write themselves in the same way as the eval shortcuts
	// record doesnt yet? have a way to distinguish b/t the literal json and the eval json, so context matters.
	if aff != affine.Record {
		ret, err = readLiteral(aff.String()+"_eval", cls, msg)
	} else {
		var obj map[string]interface{}
		if e := json.Unmarshal(msg, &obj); e != nil {
			err = e
		} else if fields, e := unmarshalFields(obj); e != nil {
			err = e
		} else {
			ret = &RecordValue{Kind: cls, Fields: fields}
		}
	}
	return
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
