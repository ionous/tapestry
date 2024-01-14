package encode

import (
	"errors"
	"fmt"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/walk"
)

type Encoder struct {
	customEncoder CustomEncoder
}

func (enc *Encoder) Customize(c CustomEncoder) *Encoder {
	enc.customEncoder = c
	return enc
}

// given a tapestry command, return its plain value.
type CustomEncoder func(enc *Encoder, cmd jsn.Marshalee) (any, error)

// turn the passed tapestry command into plain values.
func (enc *Encoder) Encode(out jsn.Marshalee) (ret any, err error) {
	src := r.ValueOf(out).Elem()
	switch t := src.Type(); t.Kind() { // ugh
	default:
		err = unknownType(t)
	case r.Struct:
		switch src.NumField() {
		case 0:
			err = unknownType(t)
		case 1:
			// slots are structs containing { Value *slot }
			// this ignores the slot and encodes the contents
			// meaning the decoder needs to know what slot to read into.
			if ptr := src.Field(0); ptr.IsValid() {
				if w := walk.Walk(ptr.Elem()); w.Next() {
					ret, err = enc.writeFlow(w.Walk())
				}
			}
		default: // flow is struct > 1 field ( Markup + something )
			ret, err = enc.writeFlow(walk.Walk(src))
		}
	// slice is a []slot or []flow
	case r.Slice:
		w := walk.Walk(src)
		elType := src.Type().Elem()
		switch elType.Kind() {
		default:
			err = unknownType(t) // print the original type
		case r.Interface:
			ret, err = enc.encodeSlots(w)
		case r.Ptr:
			ret, err = enc.encodeFlows(w)
		}
	}
	return
}

// it is at the struct level
func (enc *Encoder) writeFlow(src walk.Walker) (ret any, err error) {
	if fix, e := unpack(src); e != nil {
		err = e
	} else if res, e := enc.customEncode(fix); e == nil {
		ret = res
	} else if e != compact.Unhandled && !errors.Is(e, compact.Unhandled) {
		err = e
	} else {
		var out FlowBuilder
		if m, e := getMarkup(src); e != nil {
			err = e
		} else if fix, ok := fix.(composer.Composer); !ok {
			err = fmt.Errorf("%s is not a flow", src.Value().Type())
		} else {
			out.WriteLede(fix.Compose().GetLede())
			out.SetMarkup(m)
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

func (enc *Encoder) customEncode(cmd jsn.Marshalee) (ret any, err error) {
	if c := enc.customEncoder; c == nil {
		err = compact.Unhandled
	} else {
		ret, err = c(enc, cmd)
	}
	return
}

func unknownType(t r.Type) error {
	return fmt.Errorf("unknown type %s(%s)", t.Kind(), t.String())
}

func unpack(src walk.Walker) (ret jsn.Marshalee, err error) {
	if fix, ok := src.Value().Addr().Interface().(jsn.Marshalee); !ok {
		err = fmt.Errorf("%s is not a flow", src.Value().Type())
	} else {
		ret = fix
	}
	return
}

// read markup from the passed walker
func getMarkup(src walk.Walker) (ret map[string]any, err error) {
	if markup := src.Markup(); markup.IsValid() {
		if m, ok := markup.Interface().(map[string]any); !ok {
			err = fmt.Errorf("expected markup, have %s", markup.Type())
		} else {
			ret = m
		}
	}
	return
}
