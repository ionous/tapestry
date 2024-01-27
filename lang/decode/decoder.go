package decode

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

type Decoder struct {
	signatures    SignatureTable
	customDecoder CustomDecoder
	patterns      PatternDecoder
}

// provide translation from command signatures to actual commands.
// ex. "Tapestry:" to *story.StoryFile
func (dec *Decoder) Signatures(t ...map[uint64]typeinfo.Instance) *Decoder {
	dec.signatures = t
	return dec
}

// customize the handling of specific slots.
func (dec *Decoder) Customize(c CustomDecoder) *Decoder {
	dec.customDecoder = c
	return dec
}

// handle conversion of unknown commands to scripted function calls.
func (dec *Decoder) Patternize(p PatternDecoder) *Decoder {
	dec.patterns = p
	return dec
}

// create an arbitrary command from arbitrary data
// ex. boolean literals are stored as actual bool values.
type CustomDecoder func(dec *Decoder, slot *typeinfo.Slot, plainData any) (typeinfo.Instance, error)

// handle pattern parsing.
type PatternDecoder func(dec *Decoder, slot *typeinfo.Slot, msg compact.Message) (typeinfo.Instance, error)

// given a desired output structure, read the passed plain data
func (dec *Decoder) Decode(out typeinfo.Instance, plainData any) (err error) {
	w := inspect.Walk(out)
	if t := w.TypeInfo(); !w.Repeating() {
		if slot, ok := t.(*typeinfo.Slot); ok {
			err = dec.decodeSlot(w, slot, plainData)
		} else {
			if msg, e := ParseMessage(plainData); e != nil {
				err = e
			} else {
				err = dec.readMsg(msg, w)
			}
		}
	} else {
		if slot, ok := t.(*typeinfo.Slot); ok {
			err = dec.repeatSlot(w, slot, plainData)
		} else {
			err = dec.repeatFlow(w, plainData)
		}
	}
	return
}

// assumes that it is at the start of a flow container
func (dec *Decoder) readMsg(msg compact.Message, out inspect.It) (err error) {
	// technically, the iterator should be the source of truth here.
	// but since the args were built from the signature
	// and the signature was matched against the registry:
	// we know we have the right type and args.
Break:
	for i, cnt, it := 0, len(msg.Args), out; i < cnt; i++ {
		p := msg.Labels[i] // fix: clear fields that get skipped?
		if f, ok := nextField(&it, p); !ok {
			err = errors.New("signature mismatch")
		} else {
			arg, out := msg.Args[i], it.RawValue()
			switch t := f.Type.(type) {
			default:
				if !f.Private { // private fields dont have typeinfo, and thats okay.
					err = fmt.Errorf("unhandled type %s", t.TypeName())
				}

			case *typeinfo.Str:
				if f.Repeats {
					err = decodeStrings(out, t, arg)
				} else {
					err = decodeString(out, t, arg)
				}

			case *typeinfo.Num:
				if f.Repeats {
					err = decodeNumbers(out, arg)
				} else if ok := it.SetValue(arg); !ok {
					err = fmt.Errorf("couldnt assign from %T", arg)
				}

			case *typeinfo.Flow:
				if arg != nil {
					if f.Repeats {
						err = dec.repeatFlow(it.Walk(), arg)
					} else if msg, e := ParseMessage(arg); e != nil {
						err = e
					} else {
						err = dec.readMsg(msg, it.Walk())
					}
				}

			case *typeinfo.Slot:
				if arg != nil {
					if slot := it.Walk(); f.Repeats {
						err = dec.repeatSlot(slot, t, arg)
					} else {
						err = dec.decodeSlot(slot, t, arg)
					}
				}
			}
		}
		if err != nil {
			err = fmt.Errorf("%q(@%s:) %w", msg.Key, p, err)
			break Break
		}
	} // for
	if err == nil && len(msg.Markup) > 0 {
		m := out.Markup(true) // fix: make this a sink during decoding?
		for k, v := range msg.Markup {
			m[k] = v
		}
	}
	return
}
