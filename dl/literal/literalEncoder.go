package literal

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"

	r "reflect"
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
		if len(out.Values) == 1 {
			err = m.MarshalValue(typeName, out.Values[0])
		} else {
			err = m.MarshalValue(typeName, out.Values)
		}

	case *TextValues:
		if len(out.Values) == 1 {
			err = m.MarshalValue(typeName, out.Values[0])
		} else {
			err = m.MarshalValue(typeName, out.Values)
		}

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

func CompactSlotDecoder(m jsn.Marshaler, slot jsn.SlotBlock, msg r.Value) (err error) {
	if ptr, e := readLiteral(slot.GetType(), "", msg); e != nil {
		err = e
	} else if !slot.SetSlot(ptr) {
		err = errutil.New("unexpected error setting slot")
	}
	return
}

func ReadLiteral(aff affine.Affinity, kind string, msg r.Value) (ret LiteralValue, err error) {
	// most literals write themselves in the same way as the eval shortcuts
	// record doesnt yet? have a way to distinguish b/t the literal json and the eval json, so context matters.
	if aff != affine.Record {
		ret, err = readLiteral(aff.String()+"_eval", kind, msg)
	} else if fields, e := unmarshalFields(msg); e != nil {
		err = e
	} else {
		ret = &RecordValue{Kind: kind, Fields: fields}
	}
	return
}

func readLiteral(typeName, kind string, msg r.Value) (ret LiteralValue, err error) {
	// when decoding, we havent created the command yet ( we're doing that now )
	// so we have to switch on the typename not the value in the slot.
	switch typeName {
	default:
		// note: trying to read a record literal directly into a record eval slot wouldnt work well
		// for one, it would have to know what type the record is.
		// RecordValue ( "record:fields:<kind>, <field values>" ) is the thing for now.
		err = chart.Unhandled("CustomSlot")

	case rt.BoolEval_Type:
		if msg.Kind() != r.Bool {
			err = chart.Unhandled(typeName)
		} else {
			ret = &BoolValue{Value: msg.Bool(), Kind: kind}
		}

	case rt.NumberEval_Type:
		if msg.Kind() != r.Float64 {
			err = chart.Unhandled(typeName)
		} else {
			ret = &NumValue{Value: msg.Float(), Kind: kind}
		}

	case rt.TextEval_Type:
		switch msg.Kind() {
		case r.String:
			ret = &TextValue{Value: msg.String(), Kind: kind}
		case r.Slice:
			if lines, e := cin.SliceLines(msg); e != nil {
				err = chart.Unhandled(typeName)
			} else {
				ret = &TextValue{Value: lines, Kind: kind}
			}
		default:
			err = chart.Unhandled(typeName)
		}

	case rt.NumListEval_Type:
		switch msg.Kind() {
		case r.Slice:
			if vs, ok := cin.SliceFloats(msg); !ok {
				err = chart.Unhandled(typeName)
			} else {
				ret = &NumValues{Values: vs, Kind: kind}
			}
		case r.Float64:
			ret = &NumValues{Values: []float64{msg.Float()}, Kind: kind}
		default:
			err = chart.Unhandled(typeName)
		}

	case rt.TextListEval_Type:
		switch msg.Kind() {
		case r.Slice:
			if vs, ok := cin.SliceStrings(msg); !ok {
				err = chart.Unhandled(typeName)
			} else {
				ret = &TextValues{Values: vs, Kind: kind}
			}
		case r.String:
			ret = &TextValues{Values: []string{msg.String()}, Kind: kind}
		default:
			err = chart.Unhandled(typeName)
		}
	}
	return
}
