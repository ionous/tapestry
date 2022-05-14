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
		err = m.MarshalValue(typeName, out.Value)

	case *NumValue:
		err = m.MarshalValue(typeName, out.Value)

	case *TextValue:
		err = m.MarshalValue(typeName, out.Value)

	case *NumValues:
		err = m.MarshalValue(typeName, out.Values)

	case *TextValues:
		err = m.MarshalValue(typeName, out.Values)

	// records dont want to contain the names of their records
	// since that's generally recoverable from knowing the names of the fields
	// it might be better to allow that though, and only remove that info on the special compact encoding
	// otherwise writing are reading arent exactly idempotent ( writes from fields but generates a record )
	case *FieldList:
		obj := make(map[string]interface{})
		if e := marshalFields(obj, out.Fields); e != nil {
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

func ReadLiteral(aff affine.Affinity, kind string, msg json.RawMessage) (ret LiteralValue, err error) {
	// most literals write themselves in the same way as the eval shortcuts
	// record doesnt yet? have a way to distinguish b/t the literal json and the eval json, so context matters.
	if aff != affine.Record {
		ret, err = readLiteral(aff.String()+"_eval", kind, msg)
	} else {
		var obj map[string]interface{}
		if e := json.Unmarshal(msg, &obj); e != nil {
			err = e
		} else if fields, e := unmarshalFields(obj); e != nil {
			err = e
		} else {
			ret = &RecordValue{Kind: kind, Fields: fields}
		}
	}
	return
}

func readLiteral(typeName, kind string, msg json.RawMessage) (ret LiteralValue, err error) {
	// when decoding, we havent created the command yet ( we're doing that now )
	// so we have to switch on the typename not the value in the slot.
	switch typeName {
	default:
		err = chart.Unhandled("CustomSlot")

	case rt.BoolEval_Type:
		var val bool
		if e := json.Unmarshal(msg, &val); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			ret = &BoolValue{Value: val, Kind: kind}
		}

	case rt.NumberEval_Type:
		var val float64
		if e := json.Unmarshal(msg, &val); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			ret = &NumValue{Value: val, Kind: kind}
		}

	case rt.TextEval_Type:
		var val string
		if e := json.Unmarshal(msg, &val); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			ret = &TextValue{Value: val, Kind: kind}
		}

	case rt.NumListEval_Type:
		var val []float64
		if e := json.Unmarshal(msg, &val); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			ret = &NumValues{Values: val, Kind: kind}
		}

	case rt.TextListEval_Type:
		var val []string
		if e := json.Unmarshal(msg, &val); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			ret = &TextValues{Values: val, Kind: kind}
		}
		// note: trying to read a record literal into a record eval slot wouldnt work well
		// it could differentiate b/t a record and command --
		// and it would have to know what type the record is.
		// could potentially think of it as a constructor
		// "<recordName>:args:" -- might take some work to get there
		// MakeRecord is the thing for now...
	}
	return
}
