package decode

import (
	"git.sr.ht/~ionous/tapestry/lang/walk"
)

func (dec *Decoder) repeatFlow(out walk.Walker, val any) (err error) {
	if els, ok := val.([]any); !ok { // single values can stand in as a slice of one
		err = dec.repeatFlow(out, []any{val})
	} else {
		out.Resize(len(els))
		for i := 0; out.Next(); i++ {
			if msg, e := parseMessage(els[i]); e != nil {
				err = e
			} else {
				err = dec.readMsg(msg, out.Walk())
			}
		}
	}
	return
}

func (dec *Decoder) repeatSlot(out walk.Walker, val any) (err error) {
	if els, ok := val.([]any); !ok { // single values can stand in as a slice of one
		err = dec.repeatSlot(out, []any{val})
	} else {
		slot := walk.SlotName(out.Value().Type().Elem())
		out.Resize(len(els))
		for i := 0; out.Next(); i++ {
			if e := dec.slotData(slot, out.Value(), els[i]); e != nil {
				err = e
				break
			}
		}
	}
	return
}
