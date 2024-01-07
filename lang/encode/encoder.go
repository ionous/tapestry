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
	// okay --- now what are all the wacky overrides:
	// assign, coreEncode, literalEncoder
	CustomEncoder
}

// given an command value, return its json-like value
// fix: cmd should have an interface
type CustomEncoder func(enc *Encoder, cmd any) (any, error)

func (enc *Encoder) Marshal(m *map[string]any, ptr any) (err error) {
	var out FlowBuilder
	if e := enc.writeFlow(&out, walk.Walk(r.ValueOf(ptr).Elem())); e != nil {
		err = e
	} else {
		*m = out.FinalizeMap()
	}
	return
}

// it is at the struct level
func (enc *Encoder) writeFlow(out *FlowBuilder, it walk.Walker) (err error) {
	if fix, ok := it.Value().Addr().Interface().(composer.Composer); !ok {
		err = fmt.Errorf("%s is not a flow", it.Value().Type())
	} else {
		lede := fix.Compose().GetLede()
		out.WriteLede(lede)
		// skip the markup
		for i, cnt := 0, it.ContainerLen()-1; i < cnt; i++ {
			it.Next()
			f, val := it.Field(), it.Value()
			if !f.Optional() || !val.IsZero() {
				switch t := f.SpecType(); t {
				default:
					err = fmt.Errorf("unhandled type %s", t)

				case walk.Str:
					err = WriteEnum(out, f, val.String())

				case walk.Value:
					err = out.WriteField(f, val.Interface())

				case walk.Flow:
					if res, e := enc.customEncode(out, val); e == nil {
						err = out.WriteField(f, res)
					} else if e != compact.Unhandled && !errors.Is(e, compact.Unhandled) {
						err = e
					} else if f.Repeats() {
						if slice, e := enc.encodeFlows(it.Walk()); e != nil {
							err = e
						} else {
							err = out.WriteField(f, slice)
						}
					} else {
						var sub FlowBuilder
						if e := enc.writeFlow(&sub, it); e != nil {
							err = e
						} else {
							err = out.WriteField(f, sub.FinalizeMap())
						}
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
					} else {
						var sub FlowBuilder
						if e := enc.writeFlow(&sub, slot.Walk()); e != nil {
							err = e
						} else {
							err = out.WriteField(f, sub.FinalizeMap())
						}
					}
				}
			}
			if err != nil {
				err = fmt.Errorf("%q(@%s) %w", lede, f.Name(), err)
				break
			}
		} // (end for)
		if err == nil && it.Next() { // handle markup
			last := it.Value()
			if m, ok := last.Interface().(map[string]any); !ok {
				err = fmt.Errorf("expected markup, have %s in %T", last.Type(), fix)
			} else {
				out.SetMarkup(m)
			}
		}
	}
	return
}

func (enc *Encoder) customEncode(out *FlowBuilder, val r.Value) (ret any, err error) {
	if c := enc.CustomEncoder; c == nil {
		err = compact.Unhandled
	} else {
		ret, err = c(enc, val.Interface)
	}
	return
}
