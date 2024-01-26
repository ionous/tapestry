package literal

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"

	"git.sr.ht/~ionous/tapestry/lang/decode"
	"github.com/ionous/errutil"
)

// Write literal commands as plain values.
// ex. BoolValue becomes a bool (true or false) in the output data.
//
// note: TextValues and NumValues of containing a single value
// and serialized as that single value ( [6]-> 6 )
// because, in theory, that can help simply the specification for authors.
func CustomEncoder(enc *encode.Encoder, op typeinfo.Inspector) (ret any, err error) {
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
			ret = anySlice(out.Values)
		}
	case *TextValues:
		if len(out.Values) == 1 {
			ret = out.Values[0]
		} else {
			ret = anySlice(out.Values)
		}
	}
	return
}

// convert slices of specific types to slices of any
// the rationale here is that plain values slices and maps
// always deserialize into the "any" type, therefore
// to be idempotent marshaling should too
func anySlice[V any](els []V) []any {
	slice := make([]any, len(els))
	for i, v := range els {
		slice[i] = v
	}
	return slice
}

func CustomDecoder(_ *decode.Decoder, slot *typeinfo.Slot, body any) (typeinfo.Inspector, error) {
	return readLiteral(slot.Name, "", body)
}

func DecodeLiteral(slot *typeinfo.Slot, body any) (ret typeinfo.Inspector, err error) {
	return readLiteral(slot.Name, "", body)
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

type literalCommand interface {
	typeinfo.Inspector
	LiteralValue
}

func readLiteral(typeName, kind string, val any) (ret literalCommand, err error) {
	// when decoding, we havent created the command yet ( we're doing that now )
	// so we have to switch on the typename not the value in the slot.
	switch typeName {
	default:
		// note: trying to read a record literal directly into a record eval slot wouldnt work well
		// for one, it would have to know what type the record is.
		// RecordValue ( "record:fields:<kind>, <field values>" ) is the thing for now.
		err = compact.Unhandled

	case rtti.Zt_BoolEval.Name:
		if v, ok := val.(bool); !ok {
			err = compact.Unhandled
		} else {
			ret = &BoolValue{Value: v, Kind: kind}
		}

	case rtti.Zt_NumberEval.Name:
		if v, ok := val.(float64); !ok {
			err = compact.Unhandled
		} else {
			ret = &NumValue{Value: v, Kind: kind}
		}

	case rtti.Zt_TextEval.Name:
		switch v := val.(type) {
		case string:
			ret = &TextValue{Value: v, Kind: kind}
		case []any:
			if lines, ok := compact.SliceLines(v); !ok {
				err = compact.Unhandled
			} else {
				ret = &TextValue{Value: lines, Kind: kind}
			}
		default:
			err = compact.Unhandled
		}

	case rtti.Zt_NumListEval.Name:
		switch v := val.(type) {
		case []any:
			if vs, ok := compact.SliceFloats(v); !ok {
				err = compact.Unhandled
			} else {
				ret = &NumValues{Values: vs, Kind: kind}
			}
		case float64:
			ret = &NumValues{Values: []float64{v}, Kind: kind}
		default:
			err = compact.Unhandled
		}

	case rtti.Zt_TextListEval.Name:
		switch v := val.(type) {
		case []any:
			if vs, ok := compact.SliceStrings(v); !ok {
				err = compact.Unhandled
			} else {
				ret = &TextValues{Values: vs, Kind: kind}
			}
		case string:
			ret = &TextValues{Values: []string{v}, Kind: kind}
		default:
			err = compact.Unhandled
		}
	}
	return
}
