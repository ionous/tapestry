package shortcut

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

func Encoder(enc *encode.Encoder, op typeinfo.Instance) (ret any, err error) {
	switch op := op.(type) {
	case rtti.Address:
		if str, ok := WriteDots(op); !ok {
			err = compact.Unhandled("address")
		} else {
			ret = str
		}

	case *call.CallPattern:
		ret, err = call.CustomEncoder(enc, op)

	case *object.AtField:
		if a, ok := op.FieldName.(*literal.TextValue); ok {
			ret = a.Value
		} else {
			err = compact.Unhandled("field")
		}

	case *object.AtIndex:
		if a, ok := op.Index.(*literal.NumValue); ok {
			ret = a.Value
		} else {
			err = compact.Unhandled("field")
		}

	case *literal.TextValue:
		// if the text starts with an @, skip it:
		// ( ie. dont confuse the rare text literal starting with an ampersand, with GetVar )
		if str := op.Value; len(str) > 0 && str[0] == '@' {
			err = compact.Unhandled("text value")
		} else {
			ret, err = literal.CustomEncoder(enc, op)
		}

	default:
		ret, err = literal.CustomEncoder(enc, op)
	}
	return
}

func Decoder(dec *decode.Decoder, slot *typeinfo.Slot, body any) (ret typeinfo.Instance, err error) {
	// switching on the slot ptr's type seems like it should work, but only results in untyped interfaces
	switch slot {
	default:
		err = compact.Unhandled("shortcut decoder")
	case
		// reading from a variable:
		&rtti.Zt_BoolEval,
		&rtti.Zt_NumEval,
		&rtti.Zt_TextEval,
		&rtti.Zt_NumListEval,
		&rtti.Zt_TextListEval,
		&rtti.Zt_RecordEval,
		&rtti.Zt_RecordListEval,
		// writing to a variable:
		&rtti.Zt_Address:
		//
		if str, ok := body.(string); !ok || len(str) == 0 {
			ret, err = literal.DecodeLiteral(slot, body)
		} else {
			var clip NotShort
			if a, e := ReadDots(str); e == nil {
				ret = a.(typeinfo.Instance)
			} else if !errors.As(e, &clip) {
				err = e
			} else {
				// use decode literal, could be string in a list.
				ret, err = literal.DecodeLiteral(slot, str[clip:])
			}
		}
	case &object.Zt_Dot:
		switch v := body.(type) {
		case string:
			ret = &object.AtField{
				FieldName: literal.T(v),
			}
		case float64:
			ret = &object.AtIndex{
				Index: literal.F(v),
			}
		default:
			err = compact.Unhandled("dot decoder")
		}
	}
	return
}
