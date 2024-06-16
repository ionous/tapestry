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
func CustomEncoder(enc *encode.Encoder, op typeinfo.Instance) (ret any, err error) {
	switch out := op.(type) {
	default:
		// record: marshalFields?
		err = compact.Unhandled("custom literal")
	case *BoolValue:
		ret = out.Value
	case *NumValue:
		ret = out.Value
	case *TextValue:
		ret = out.Value
	case *NumList:
		if len(out.Values) == 1 {
			ret = out.Values[0]
		} else {
			ret = anySlice(out.Values)
		}
	case *TextList:
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

func CustomDecoder(_ *decode.Decoder, slot *typeinfo.Slot, body any) (typeinfo.Instance, error) {
	return readLiteral(slot.Name, "", body)
}

func DecodeLiteral(slot *typeinfo.Slot, body any) (ret typeinfo.Instance, err error) {
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
		ret = &RecordValue{KindName: kind, Fields: fields}
	}
	return
}

type literalCommand interface {
	typeinfo.Instance
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
		err = compact.Unhandled("literal value")

	case rtti.Zt_BoolEval.Name:
		if v, ok := val.(bool); !ok {
			err = compact.Unhandled("bool")
		} else {
			ret = &BoolValue{Value: v}
		}

	case rtti.Zt_NumEval.Name:
		if v, ok := val.(float64); !ok {
			err = compact.Unhandled("float")
		} else {
			ret = &NumValue{Value: v}
		}

	case rtti.Zt_TextEval.Name:
		switch v := val.(type) {
		case string:
			ret = &TextValue{Value: v, KindName: kind}
		case []any:
			if lines, ok := compact.JoinLines(v); !ok {
				err = compact.Unhandled("lines")
			} else {
				ret = &TextValue{Value: lines, KindName: kind}
			}
		default:
			err = compact.Unhandled("text")
		}

	case rtti.Zt_NumListEval.Name:
		switch v := val.(type) {
		case []any:
			if vs, ok := compact.SliceFloats(v); !ok {
				err = compact.Unhandled("floats")
			} else {
				ret = &NumList{Values: vs}
			}
		case float64:
			ret = &NumList{Values: []float64{v}}
		default:
			err = compact.Unhandled("numbers")
		}

	case rtti.Zt_TextListEval.Name:
		switch v := val.(type) {
		case []any:
			if vs, ok := compact.SliceStrings(v); !ok {
				err = compact.Unhandled("strings")
			} else {
				ret = &TextList{Values: vs, KindName: kind}
			}
		case string:
			ret = &TextList{Values: []string{v}, KindName: kind}
		default:
			err = compact.Unhandled("text values")
		}
	}
	return
}
