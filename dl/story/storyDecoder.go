package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/rt"
)

func Decode(msg map[string]any) (ret StoryFile, err error) {
	err = DecodeMessage(&ret, msg)
	return
}

// Create a story dl from native maps and slices.
func DecodeMessage(ptr composer.Composer, msg map[string]any) error {
	d := decode.Decoder{
		SignatureTable: AllSignatures,
		CustomDecoder:  core.CoreDecoder,
		PatternDecoder: TryPattern,
	}
	return d.Unmarshal(ptr, msg)
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
			var val rt.Assignment
			if e := dec.UnmarshalSlot(&val, args[i]); e != nil {
				err = e
				break
			} else {
				ret = append(ret, assign.Arg{Name: p.Label, Value: val})
			}
		}
	}
	return
}
