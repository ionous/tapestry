package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
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
		Patterns(TryPattern)
	return dec.Decode(out, msg)
}

func TryPattern(dec *decode.Decoder, slot string, msg compact.Message) (ret any, err error) {
	switch slot {
	default:
		err = compact.Unhandled
	case StoryStatement_Type:
		if args, e := tryPatternArgs(dec, slot, msg); e != nil {
			err = e
		} else {
			ret = &CallMacro{MacroName: msg.Name, Arguments: args}
		}
	case
		rt.Execute_Type,
		rt.BoolEval_Type,
		rt.NumberEval_Type,
		rt.TextEval_Type,
		rt.RecordEval_Type,
		rt.NumListEval_Type,
		rt.TextListEval_Type,
		rt.RecordListEval_Type:
		if args, e := tryPatternArgs(dec, slot, msg); e != nil {
			err = e
		} else {
			ret = &assign.CallPattern{PatternName: msg.Name, Arguments: args}
		}
	}
	return
}

func tryPatternArgs(dec *decode.Decoder, slot string, msg compact.Message) (ret []assign.Arg, err error) {
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
