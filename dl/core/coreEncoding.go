package core

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
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

	case *assign.VariableDot:
		// write variables as a string prepended by @
		// fix: it'd be nice if all parts were literals to write dot/bracket syntax a.b[5]
		// fix: it'd be nicest if this could use package express to handle the parsing.
		if name, ok := op.Name.(*literal.TextValue); len(op.Dot) == 0 && ok {
			ret = "@" + name.Value
		} else {
			err = compact.Unhandled("variable ref")
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
		if str := getVariableString(body); len(str) > 0 {
			ret = assign.Variable(str)
		} else {
			ret, err = literal.DecodeLiteral(slot, body)
		}
	}
	return
}

func getVariableString(val any) (ret string) {
	if str, ok := val.(string); ok && len(str) > 0 && str[0] == '@' {
		ret = str[1:]
	}
	return
}
