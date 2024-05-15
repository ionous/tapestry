package encode

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

type Encoder struct {
	customEncoder CustomEncoder
}

func (enc *Encoder) Customize(c CustomEncoder) *Encoder {
	enc.customEncoder = c
	return enc
}

// given a tapestry flow, return its plain value.
type CustomEncoder func(*Encoder, typeinfo.Instance) (any, error)

// turn the passed tapestry command into plain values.
func (enc *Encoder) Encode(i typeinfo.Instance) (ret any, err error) {
	w := inspect.Walk(i)
	if t := w.TypeInfo(); !w.Repeating() {
		if _, ok := t.(*typeinfo.Slot); !ok {
			ret, err = enc.writeFlow(w)
		} else {
			// slots are structs containing { Value slot }
			// this ignores the slot and encodes the contents
			// meaning the decoder needs to know what slot to read into.
			// fix? add a GetSlot/SetSlot to inspect or typeinfo?
			if w.Next() {
				ret, err = enc.writeFlow(w.Walk())
			}
		}
	} else {
		if _, ok := t.(*typeinfo.Slot); !ok {
			ret, err = enc.encodeFlows(w)
		} else {
			ret, err = enc.encodeSlots(w)
		}
	}
	return
}

// it is at the struct level
func (enc *Encoder) writeFlow(src inspect.It) (ret any, err error) {
	if fix, e := unpack(src); e != nil {
		err = e
	} else if res, e := enc.customEncode(fix); e == nil {
		ret = res
	} else if !compact.IsUnhandled(e) {
		err = e
	} else {
		var out FlowBuilder
		flow := src.TypeInfo().(*typeinfo.Flow)
		out.WriteLede(flow.Lede)
		out.SetMarkup(src.Markup(false))
		for it := src; it.Next(); { // -1 to end before the markup
			f, val := it.Term(), it.RawValue()
			if !f.Private && (!f.Optional || !val.IsZero()) {
				switch t := f.Type.(type) {
				default:
					err = fmt.Errorf("unhandled type %s", t.TypeName())

				case *typeinfo.Str, *typeinfo.Num:
					if f.Repeats {
						vs := encodeValues(it.Walk())
						out.WriteArg(f.Label, vs)
					} else {
						v := it.CompactValue()
						out.WriteArg(f.Label, v)
					}

				case *typeinfo.Flow:
					if f.Repeats {
						if slice, e := enc.encodeFlows(it.Walk()); e != nil {
							err = e
						} else {
							out.WriteArg(f.Label, slice)
						}
					} else if res, e := enc.writeFlow(it); e != nil {
						err = e
					} else {
						out.WriteArg(f.Label, res)
					}

				case *typeinfo.Slot:
					if f.Repeats {
						if slice, e := enc.encodeSlots(it.Walk()); e != nil {
							err = e
						} else {
							out.WriteArg(f.Label, slice)
						}
					} else if slot := it.Walk(); !slot.Next() {
						// ^ walk into the field to get the container of the slot
						// ^ next, check for empty contents; write nil if so.
						out.WriteArg(f.Label, nil)
					} else if res, e := enc.writeFlow(slot.Walk()); e != nil {
						err = e
					} else {
						out.WriteArg(f.Label, res)
					}
				}
			}
			if err != nil {
				err = fmt.Errorf("%T (@%s) %w", fix, f.Name, err)
				break
			}
		} // (end for)
		if err == nil {
			ret = out.FinalizeMap()
		}
	}
	return
}

func (enc *Encoder) customEncode(cmd typeinfo.Instance) (ret any, err error) {
	if c := enc.customEncoder; c == nil {
		err = compact.Unhandled("custom encoder")
	} else {
		ret, err = c(enc, cmd)
	}
	return
}

func unpack(src inspect.It) (ret typeinfo.Instance, err error) {
	v := src.RawValue()
	if fix, ok := v.Addr().Interface().(typeinfo.Instance); !ok {
		err = fmt.Errorf("%s is not a flow", v.Type())
	} else {
		ret = fix
	}
	return
}
