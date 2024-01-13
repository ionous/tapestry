package decode

import (
	"errors"
	"fmt"
	r "reflect"
	"unicode"

	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/walk"
)

type Decoder struct {
	SignatureTable
	CustomDecoder
	PatternDecoder
}

// create an arbitrary command from arbitrary data
// ex. boolean literals are stored as actual bool values.
// fix: return should be composer
type CustomDecoder func(dec *Decoder, slot string, body any) (any, error)

// handle pattern parsing.
// fix: return should be composer
type PatternDecoder func(dec *Decoder, slot string, msg compact.Message) (any, error)

func (dec *Decoder) Unmarshal(ptr composer.Composer, m map[string]any) (err error) {
	if msg, e := parseMessage(m); e != nil {
		err = e
	} else {
		err = dec.readMsg(msg, walk.Walk(r.ValueOf(ptr).Elem()))
	}
	return
}

// fix: slot should be an interface of some useful sort
// currently: address of a slot.
func (dec *Decoder) UnmarshalSlot(slotptr, data any) (err error) {
	out := r.ValueOf(slotptr).Elem()
	slot := walk.SlotName(out.Type())
	return dec.slotData(slot, out, data)
}

func (dec *Decoder) UnmarshalSlice(ptr r.Value, msgs []any) (err error) {
	panic("xx")
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
						err = dec.repeatSlot(it.Walk(), arg)
					} else {
						slot := walk.SlotName(out.Type())
						err = dec.slotData(slot, out, arg)
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
