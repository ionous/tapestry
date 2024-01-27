package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

func Encode(file *StoryFile) (ret any, err error) {
	var enc encode.Encoder
	return enc.
		Customize(core.CustomEncoder).
		Encode(file)
}

func Decode(out *StoryFile, msg map[string]any) error {
	return DecodeMessage(out, msg)
}

// Create a story dl from native maps and slices.
func DecodeMessage(out typeinfo.Inspector, msg map[string]any) error {
	var dec decode.Decoder
	dec.Signatures(AllSignatures...).
		Customize(core.CustomDecoder).
		Patternize(DecodePattern)
	return dec.Decode(out, msg)
}

func DecodePattern(dec *decode.Decoder, slot *typeinfo.Slot, msg compact.Message) (ret typeinfo.Inspector, err error) {
	switch slot {
	default:
		err = compact.Unhandled("pattern")
	case &Zt_StoryStatement:
		if args, e := tryPatternArgs(dec, msg); e != nil {
			err = e
		} else {
			ret = &CallMacro{MacroName: msg.Name, Arguments: args}
		}
	case
		&rtti.Zt_Execute,
		&rtti.Zt_BoolEval,
		&rtti.Zt_NumberEval,
		&rtti.Zt_TextEval,
		&rtti.Zt_RecordEval,
		&rtti.Zt_NumListEval,
		&rtti.Zt_TextListEval,
		&rtti.Zt_RecordListEval:
		if args, e := tryPatternArgs(dec, msg); e != nil {
			err = e
		} else {
			ret = &assign.CallPattern{PatternName: msg.Name, Arguments: args}
		}
	}
	return
}

func tryPatternArgs(dec *decode.Decoder, msg compact.Message) (ret []assign.Arg, err error) {
	for i, p := range msg.Labels {
		var val rtti.Assignment_Slot
		if e := dec.Decode(&val, msg.Args[i]); e != nil {
			err = e
			break
		} else {
			ret = append(ret, assign.Arg{Name: p, Value: val.Value})
		}
	}
	return
}
