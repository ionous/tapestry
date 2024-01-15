package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/encode"
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
func DecodeMessage(ptr jsn.Marshalee, msg map[string]any) error {
	var dec decode.Decoder
	dec.Signatures(AllSignatures...).
		Customize(core.CustomDecoder).
		Patterns(TryPattern)
	return dec.Decode(ptr, msg)
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
	if args, e := msg.Args(); e != nil {
		err = e
	} else {
		for i, p := range msg.Params {
			var val rt.Assignment_Slot
			if e := dec.Decode(&val, args[i]); e != nil {
				err = e
				break
			} else {
				ret = append(ret, assign.Arg{Name: p.Label, Value: val.Value})
			}
		}
	}
	return
}
