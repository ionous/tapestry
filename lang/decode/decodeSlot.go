package decode

import (
	"errors"
	"fmt"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/walk"
)

func (dec *Decoder) slotData(slot string, out r.Value, data any) (err error) {
	switch e := dec.customDecode(slot, out, data); {
	default:
		err = e
	case e == compact.Unhandled || errors.Is(e, compact.Unhandled):
		if msg, e := parseMessage(data); e != nil {
			err = e
		} else {
			switch rptr, ok := dec.SignatureTable.Create(slot, msg.Key); {
			default:
				err = fmt.Errorf("unknown signature slot=%q flow=%q", slot, msg.Key)

			case ok:
				dst := walk.MakeWalker(rptr.Elem())
				if e := dec.readMsg(msg, dst); e != nil {
					err = e
				} else {
					out.Set(rptr)
				}

			case dec.PatternDecoder != nil:
				if ptr, e := dec.PatternDecoder(dec, slot, msg); e != nil {
					err = e
				} else if rptr := r.ValueOf(ptr); !assignSlot(out, rptr) {
					err = fmt.Errorf("couldnt assign decoded pattern out:%s in:%s", out.Type(), rptr.Type())
				}
			}
		}
	}
	return
}

func (dec *Decoder) customDecode(slot string, out r.Value, arg any) (err error) {
	if dec.CustomDecoder == nil {
		err = compact.Unhandled
	} else if ptr, e := dec.CustomDecoder(dec, slot, arg); e != nil {
		err = e
	} else if rptr := r.ValueOf(ptr); !assignSlot(out, rptr) {
		err = fmt.Errorf("couldnt assign value of %s to slot of %s", rptr.Type(), out.Type())
	}
	return
}

func assignSlot(out, rptr r.Value) (okay bool) {
	if rptr.Type().AssignableTo(out.Type()) {
		out.Set(rptr)
		okay = true
	}
	return
}
