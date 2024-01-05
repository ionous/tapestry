package decode

import (
	"errors"
	"fmt"
	r "reflect"

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
		err = dec.readMsg(msg, walk.MakeWalker(r.ValueOf(ptr).Elem()))
	}
	return
}

// fix: slot should be an interface of some useful sort
// currently: address of a slot.
func (dec *Decoder) UnmarshalSlot(slotptr, data any) (err error) {
	out := r.ValueOf(slotptr).Elem()
	slot := slotName(out.Type())
	return dec.slotData(slot, out, data)
}

func (dec *Decoder) UnmarshalSlice(ptr r.Value, msgs []any) (err error) {
	panic("xx")
}

// assumes that it is at the start of a flow container
func (dec *Decoder) readMsg(msg compact.Message, it walk.Walker) (err error) {
	if args, e := msg.Args(); e != nil {
		err = e
	} else {
		// technically, the iterator should be the source of truth here.
		// but since the args were built from the signature
		// and the signature was matched against the registry:
		// we know we have the right type and args.
		for i, cnt := 0, len(args); i < cnt; i++ {
			p := msg.Params[i] // fix: clear fields that get skipped?
			if f, ok := nextField(&it, p); !ok {
				err = errors.New("signature mismatch")
			} else {
				arg, out := args[i], it.Value()
				switch t := it.Type(); t {
				default:
					err = fmt.Errorf("unhandled type %s", t)

				case walk.Str:
					err = SetString(out, f.Type(), arg) // fix: are there repeating strings?

				case walk.Value:
					if f.Repeats() {
						err = SetValues(out, arg)
					} else {
						err = SetValue(out, arg)
					}

				case walk.Flow:
					if f.Repeats() {
						err = dec.repeatFlow(out, arg)
					} else {
						if msg, e := parseMessage(arg); e != nil {
							err = e
						} else if subit, ok := it.Descend(); ok {
							err = dec.readMsg(msg, subit)
						}
					}

				case walk.Slot:
					if f.Repeats() {
						slot := slotName(out.Type().Elem())
						err = dec.repeatSlot(slot, out, arg)
					} else {
						slot := slotName(out.Type())
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
			// future: remove reflection, and provide access to Markup through Compose()?
			// ( would be nice b/c we could do this before everything else, w/o err testing )
			markup := findLast(it) // cause this is kind of hacky.
			if v := markup.Value(); v.Kind() == r.Map {
				v.Set(r.ValueOf(msg.Markup))
			}
		}

	}
	return
}
