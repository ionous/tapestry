package literal

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/encode"
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

func CustomEncoder(enc *encode.Encoder, op any) (ret any, err error) {
	switch out := op.(type) {
	default:
		err = compact.Unhandled

	case *BoolValue:
		ret = out.Value

	case *NumValue:
		ret = out.Value

	case *TextValue:
		ret = out.Value

	case *NumValues:
		if len(out.Values) == 1 {
			ret = out.Values[0]
		} else {
			ret = out.Values
		}

	case *TextValues:
		if len(out.Values) == 1 {
			ret = out.Values[0]
		} else {
			ret = out.Values
		}
	}
	return
}

func CompactSlotDecoder(m jsn.Marshaler, slot jsn.SlotBlock, body any) (err error) {
	if ptr, e := readLiteral(slot.GetType(), "", body); e != nil {
		err = e
	} else if !slot.SetSlot(ptr) {
		err = errutil.New("unexpected error setting slot")
	}
	return
}

func ReadLiteral(aff affine.Affinity, kind string, val any) (ret LiteralValue, err error) {
	// most literals write themselves in the same way as the eval shortcuts
	// record doesnt yet? have a way to distinguish b/t the literal json and the eval json, so context matters.
	if aff != affine.Record {
		ret, err = readLiteral(aff.String()+"_eval", kind, val)
	} else if msg, ok := val.(map[string]any); !ok {
		err = errutil.Fmt("expected a map, have %T", val)
	} else if fields, e := unmarshalFields(msg); e != nil {
		err = e
	} else {
		ret = &RecordValue{Kind: kind, Fields: fields}
	}
	return
}

func LiteralDecoder(dec *decode.Decoder, slot string, body any) (ret any, err error) {
	var u chart.Unhandled
	if ret, err = readLiteral(slot, "", body); errors.As(err, &u) {
		err = compact.Unhandled
	}
	return
}

func readLiteral(typeName, kind string, val any) (ret LiteralValue, err error) {
	// when decoding, we havent created the command yet ( we're doing that now )
	// so we have to switch on the typename not the value in the slot.
	switch typeName {
	default:
		// note: trying to read a record literal directly into a record eval slot wouldnt work well
		// for one, it would have to know what type the record is.
		// RecordValue ( "record:fields:<kind>, <field values>" ) is the thing for now.
		err = chart.Unhandled("CustomSlot")

	case rt.BoolEval_Type:
		if v, ok := val.(bool); !ok {
			err = chart.Unhandled(typeName)
		} else {
			ret = &BoolValue{Value: v, Kind: kind}
		}

	case rt.NumberEval_Type:
		if v, ok := val.(float64); !ok {
			err = chart.Unhandled(typeName)
		} else {
			ret = &NumValue{Value: v, Kind: kind}
		}

	case rt.TextEval_Type:
		switch v := val.(type) {
		case string:
			ret = &TextValue{Value: v, Kind: kind}
		case []any:
			if lines, ok := compact.SliceLines(v); !ok {
				err = chart.Unhandled(typeName)
			} else {
				ret = &TextValue{Value: lines, Kind: kind}
			}
		default:
			err = chart.Unhandled(typeName)
		}

	case rt.NumListEval_Type:
		switch v := val.(type) {
		case []any:
			if vs, ok := compact.SliceFloats(v); !ok {
				err = chart.Unhandled(typeName)
			} else {
				ret = &NumValues{Values: vs, Kind: kind}
			}
		case float64:
			ret = &NumValues{Values: []float64{v}, Kind: kind}
		default:
			err = chart.Unhandled(typeName)
		}

	case rt.TextListEval_Type:
		switch v := val.(type) {
		case []any:
			if vs, ok := compact.SliceStrings(v); !ok {
				err = chart.Unhandled(typeName)
			} else {
				ret = &TextValues{Values: vs, Kind: kind}
			}
		case string:
			ret = &TextValues{Values: []string{v}, Kind: kind}
		default:
			err = chart.Unhandled(typeName)
		}
	}
	return
}
