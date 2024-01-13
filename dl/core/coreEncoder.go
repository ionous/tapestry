package core

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/rt"
)

// fix: move some part of this into package assign?
// especially because the opposite is handled in... story?
func CustomEncoder(enc *encode.Encoder, op any) (ret any, err error) {
	switch op := op.(type) {
	case *assign.CallPattern:
		ret, err = assign.CustomEncoder(enc, op)

	case *assign.VariableRef:
		// write variables as a string prepended by @
		// fix: it'd be nice if all parts were literals to write dot/bracket syntax a.b[5]
		// fix: it'd be nicest if this could use package express to handle the parsing.
		if name, ok := op.Name.(*literal.TextValue); len(op.Dot) == 0 && ok {
			ret = "@" + name.Value
		} else {
			err = compact.Unhandled
		}

	case *literal.TextValue:
		// if the text starts with an @, skip it:
		// ( ie. dont confuse the rare text literal starting with an ampersand, with GetVar )
		if str := op.Value; len(str) > 0 && str[0] == '@' {
			err = compact.Unhandled
		} else {
			ret, err = literal.CustomEncoder(enc, op)
		}

	default:
		ret, err = literal.CustomEncoder(enc, op)
	}
	return
}

func CustomDecoder(dec *decode.Decoder, slot string, body any) (ret any, err error) {
	// switching on the slot ptr's type seems like it should work, but only results in untyped interfaces
	switch slot {
	default:
		err = compact.Unhandled
	case
		// reading from a variable:
		rt.BoolEval_Type,
		rt.NumberEval_Type,
		rt.TextEval_Type,
		rt.NumListEval_Type,
		rt.TextListEval_Type,
		rt.RecordEval_Type,
		rt.RecordListEval_Type,
		// writing to a variable:
		assign.Address_Type:
		if str := getVariableString(body); len(str) > 0 {
			ret = Variable(str)
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
