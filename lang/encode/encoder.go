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

// turn the passed tapestry data into json friendly values.
// currently, expects ptr to a tapestry command, and returns a map[string]any
// future: allow ptr to be any container type
func (enc *Encoder) Marshal(ptr any) (res any, err error) {
	return enc.writeFlow(walk.Walk(r.ValueOf(ptr).Elem()))
}

// it is at the struct level
func (enc *Encoder) writeFlow(it walk.Walker) (ret any, err error) {
	if fix, ok := it.Value().Addr().Interface().(composer.Composer); !ok {
		err = fmt.Errorf("%s is not a flow", it.Value().Type())
	} else if res, e := enc.customEncode(fix); e == nil {
		ret = res
	} else if e != compact.Unhandled && !errors.Is(e, compact.Unhandled) {
		err = e
	} else {
		var out FlowBuilder
		out.WriteLede(fix.Compose().GetLede())
		for i, cnt := 0, it.Len()-1; i < cnt; i++ { // -1 to end before the markup
			it.Next()
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
		if err == nil { // handle markup
			if e := writeMarkup(&out, it); e != nil {
				err = e
			} else {
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

// it, the element before the markup
func writeMarkup(out *FlowBuilder, it walk.Walker) (err error) {
	if it.Next() {
		last := it.Value()
		if m, ok := last.Interface().(map[string]any); !ok {
			err = fmt.Errorf("expected markup, have %s", last.Type())
		} else {
			out.SetMarkup(m)
		}
	}
	return
}
