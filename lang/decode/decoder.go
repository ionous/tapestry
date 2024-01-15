package decode

import (
	"errors"
	"fmt"
	r "reflect"
	"unicode"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/walk"
)

type Decoder struct {
	signatures    SignatureTable
	customDecoder CustomDecoder
	patterns      PatternDecoder
}

// provide translation from command signatures to actual commands.
// ex. "Tapestry:" to *story.StoryFile
func (dec *Decoder) Signatures(t ...map[uint64]any) *Decoder {
	dec.signatures = t
	return dec
}

// customize the handling of specific slots.
func (dec *Decoder) Customize(c CustomDecoder) *Decoder {
	dec.customDecoder = c
	return dec
}

// handle conversion of unknown commands to scripted function calls.
func (dec *Decoder) Patterns(p PatternDecoder) *Decoder {
	dec.patterns = p
	return dec
}

// create an arbitrary command from arbitrary data
// ex. boolean literals are stored as actual bool values.
// fix: return should be composer
type CustomDecoder func(dec *Decoder, slot string, plainData any) (any, error)

// handle pattern parsing.
// fix: return should be composer or marshalee
type PatternDecoder func(dec *Decoder, slot string, msg compact.Message) (any, error)

// given a desired output structure, read the passed plain data
func (dec *Decoder) Decode(out jsn.Marshalee, plainData any) (err error) {
	tgt := r.ValueOf(out).Elem()       // the element under the interface
	switch t := tgt.Type(); t.Kind() { // ugh
	default:
		err = unknownType(t)
	case r.Struct:
		if slot, ok := out.(jsn.SlotBlock); ok {
			err = dec.decodeSlot(tgt, plainData, slot.GetType())
		} else if msg, e := parseMessage(plainData); e != nil {
			err = e
		} else {
			err = dec.readMsg(msg, walk.Walk(tgt))
		}
	// slice is a []slot or []flow
	case r.Slice:
		w := walk.Walk(tgt)
		elType := tgt.Type().Elem()
		switch elType.Kind() {
		default:
			err = unknownType(t) // print the original type
		case r.Interface:
			typeName := out.(jsn.SliceBlock).GetType()
			err = dec.repeatSlot(w, plainData, typeName)
		case r.Struct:
			err = dec.repeatFlow(w, plainData)
		}
	}
	return
}

// assumes that it is at the start of a flow container
func (dec *Decoder) readMsg(msg compact.Message, out walk.Walker) (err error) {
	if args, e := msg.Args(); e != nil {
		err = e
	} else {
		// technically, the iterator should be the source of truth here.
		// but since the args were built from the signature
		// and the signature was matched against the registry:
		// we know we have the right type and args.
		for i, cnt, it := 0, len(args), out; i < cnt; i++ {
			p := msg.Params[i] // fix: clear fields that get skipped?
			if f, ok := nextField(&it, p); !ok {
				err = errors.New("signature mismatch")
			} else {
				arg, out := args[i], it.Value()
				switch t := f.SpecType(); t {
				default:
					err = fmt.Errorf("unhandled type %s", t)

				case walk.Str:
					err = SetString(out, f.Type(), arg) // fix: are there repeating strings?

				case walk.Value:
					if f.Repeats() {
						err = SetValues(out, arg)
					} else if ok := it.SetValue(arg); !ok {
						err = fmt.Errorf("couldnt assign from %T", arg)
					}

				case walk.Swap: // ugh.
					if f.Repeats() {
						panic("not implemented")
					} else {
						choice := newStringKey(p.Choice) // "one_to_many" -> "$ONE_TO_MANY"
						i := out.Addr().Interface()      //
						swap := i.(interface{ SetSwap(string) bool })
						if !swap.SetSwap(choice) {
							err = fmt.Errorf("swap has unexpected choice %q", choice)
						} else {
							if msg, e := parseMessage(arg); e != nil {
								err = e
							} else {
								swap := it.Walk()
								swap.Next()
								err = dec.readMsg(msg, swap.Walk())
							}
						}
					}

				case walk.Flow:
					if f.Repeats() {
						err = dec.repeatFlow(it.Walk(), arg)
					} else {
						if msg, e := parseMessage(arg); e != nil {
							err = e
						} else {
							err = dec.readMsg(msg, it.Walk())
						}
					}

				case walk.Slot:
					if f.Repeats() {
						err = dec.repeatSlot(it.Walk(), arg, nameOf(out.Type().Elem()))
					} else {
						err = dec.decodeSlot(out, arg, nameOf(out.Type()))
					}
				}
			}
			if err != nil {
				err = ParamError(msg, p, err)
				break
			}
		} // for
		if err == nil {
			if markup := out.Markup(); markup.IsValid() {
				markup.Set(r.ValueOf(msg.Markup))
			}
		}
	}
	return
}

func newStringKey(s string) string {
	rs := make([]rune, 0, len(s)+1)
	rs = append(rs, '$')
	for _, r := range s {
		rs = append(rs, unicode.ToUpper(r))
	}
	return string(rs)
}
