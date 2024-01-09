package encode

import (
	"errors"
	"fmt"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/walk"
)

type Encoder struct {
	CustomEncoder
}

// given a tapestry command, return its json-like value
// fix: cmd should have an interface
type CustomEncoder func(enc *Encoder, cmd any) (any, error)

// turn the passed tapestry command into json friendly values.
func (enc *Encoder) MarshalFlow(ptr any) (ret any, err error) {
	return enc.writeFlow(walk.Walk(r.ValueOf(ptr).Elem()))
}

// turn the passed tapestry slot into json friendly values.
func (enc *Encoder) MarshalSlot(ptr any) (ret any, err error) {
	slot := walk.Walk(r.ValueOf(ptr).Elem())
	if slot.Next() {
		ret, err = enc.writeFlow(slot.Walk())
	}
	return
}

// it is at the struct level
func (enc *Encoder) writeFlow(src walk.Walker) (ret any, err error) {
	if fix, ok := src.Value().Addr().Interface().(composer.Composer); !ok {
		err = fmt.Errorf("%s is not a flow", src.Value().Type())
	} else if res, e := enc.customEncode(fix); e == nil {
		ret = res
	} else if e != compact.Unhandled && !errors.Is(e, compact.Unhandled) {
		err = e
	} else {
		var out FlowBuilder
		if e := writeMarkup(&out, src); e != nil {
			err = e
		} else {
			out.WriteLede(fix.Compose().GetLede())
			for it := src; it.Next(); { // -1 to end before the markup
				f, val := it.Field(), it.Value()
				if !f.Optional() || !val.IsZero() {
					switch t := f.SpecType(); t {
					default:
						err = fmt.Errorf("unhandled type %s", t)

					case walk.Str:
						err = WriteEnum(&out, f, val.String())

					case walk.Value:
						err = out.WriteField(f, val.Interface())

					case walk.Flow:
						if f.Repeats() {
							if slice, e := enc.encodeFlows(it.Walk()); e != nil {
								err = e
							} else {
								err = out.WriteField(f, slice)
							}
						} else if res, e := enc.writeFlow(it); e != nil {
							err = e
						} else {
							err = out.WriteField(f, res)
						}

					case walk.Slot:
						if f.Repeats() {
							if slice, e := enc.encodeSlots(it.Walk()); e != nil {
								err = e
							} else {
								err = out.WriteField(f, slice)
							}
						} else if slot := it.Walk(); !slot.Next() {
							// ^ walk into the field to get the container of the slot
							// ^ next, check for empty contents; write nil if so.
							err = out.WriteField(f, nil)
						} else if res, e := enc.writeFlow(slot.Walk()); e != nil {
							err = e
						} else {
							err = out.WriteField(f, res)
						}
					}
				}
				if err != nil {
					err = fmt.Errorf("%T (@%s) %w", fix, f.Name(), err)
					break
				}
			} // (end for)
			if err == nil {
				ret = out.FinalizeMap()
			}
		}
	}
	return
}

func (enc *Encoder) customEncode(cmd any) (ret any, err error) {
	if c := enc.CustomEncoder; c == nil {
		err = compact.Unhandled
	} else {
		ret, err = c(enc, cmd)
	}
	return
}

func writeMarkup(out *FlowBuilder, src walk.Walker) (err error) {
	if markup := src.Markup(); markup.IsValid() {
		if m, ok := markup.Interface().(map[string]any); !ok {
			err = fmt.Errorf("expected markup, have %s", markup.Type())
		} else {
			out.SetMarkup(m)
		}
	}
	return
}
