package core

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/assign/shortcut"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// fix: move some part of this into package assign?
// especially because the opposite is handled in... story?
func CustomEncoder(enc *encode.Encoder, op typeinfo.Instance) (ret any, err error) {
	switch op := op.(type) {
	case *assign.CallPattern:
		ret, err = assign.CustomEncoder(enc, op)

	case assign.Address:
		if str, ok := shortcut.WriteDots(op); !ok {
			err = compact.Unhandled("address")
		} else {
			ret = str
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

func CustomDecoder(dec *decode.Decoder, slot *typeinfo.Slot, body any) (ret typeinfo.Instance, err error) {
	// switching on the slot ptr's type seems like it should work, but only results in untyped interfaces
	switch slot {
	default:
		err = compact.Unhandled("core decoder")
	case
		// reading from a variable:
		&rtti.Zt_BoolEval,
		&rtti.Zt_NumberEval,
		&rtti.Zt_TextEval,
		&rtti.Zt_NumListEval,
		&rtti.Zt_TextListEval,
		&rtti.Zt_RecordEval,
		&rtti.Zt_RecordListEval,
		// writing to a variable:
		&assign.Zt_Address:
		//
		if str, ok := body.(string); !ok || len(str) == 0 {
			ret, err = literal.DecodeLiteral(slot, body)
		} else {
			var clip shortcut.NotShort
			if a, e := shortcut.ReadDots(str); e == nil {
				ret = a.(typeinfo.Instance)
			} else if !errors.As(e, &clip) {
				err = e
			} else {
				// use decode literal, could be string in a list.
				ret, err = literal.DecodeLiteral(slot, str[clip:])
			}
		}
	}
	return
}
