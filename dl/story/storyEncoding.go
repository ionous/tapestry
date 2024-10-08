package story

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/lang/shortcut"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

func Encode(file *StoryFile) (ret any, err error) {
	var enc encode.Encoder
	return enc.
		Customize(shortcut.Encoder).
		Encode(file)
}

func Decode(out *StoryFile, msg map[string]any) error {
	return DecodeMessage(out, msg)
}

// Create a story dl from native maps and slices.
func DecodeMessage(out typeinfo.Instance, msg map[string]any) error {
	dec := NewDecoder()
	return dec.Decode(out, msg)
}

func NewDecoder() *decode.Decoder {
	var dec decode.Decoder
	dec.Signatures(AllSignatures...).
		Customize(shortcut.Decoder).
		Patternize(DecodePattern)
	return &dec
}

// fix: move to call ( that's where encoding lives )
func DecodePattern(dec *decode.Decoder, slot *typeinfo.Slot, msg compact.Message) (ret typeinfo.Instance, err error) {
	switch slot {
	default:
		str := fmt.Sprintf("pattern named %q", msg.Key)
		err = compact.MessageError(msg, compact.Unhandled(str))
	case
		&rtti.Zt_Execute,
		&rtti.Zt_BoolEval,
		&rtti.Zt_NumEval,
		&rtti.Zt_TextEval,
		&rtti.Zt_RecordEval,
		&rtti.Zt_NumListEval,
		&rtti.Zt_TextListEval,
		&rtti.Zt_RecordListEval:
		if args, e := tryPatternArgs(dec, msg); e != nil {
			err = e
		} else {
			ret = &call.CallPattern{PatternName: msg.Lede, Arguments: args}
		}
	}
	return
}

func tryPatternArgs(dec *decode.Decoder, msg compact.Message) (ret []call.Arg, err error) {
	for i, p := range msg.Labels {
		var val rtti.Assignment_Slot
		if e := dec.Decode(&val, msg.Args[i]); e != nil {
			err = compact.MessageError(msg, e)
			break
		} else {
			ret = append(ret, call.Arg{Name: p, Value: val.Value})
		}
	}
	return
}
