package decode

import (
	r "reflect"

	"git.sr.ht/~ionous/tapestry/lang/walk"
)

func (dec *Decoder) repeatFlow(out r.Value, val any) (err error) {
	if els, ok := val.([]any); !ok { // single values can stand in as a slice of one
		err = dec.repeatFlow(out, []any{val})
	} else if cnt := len(els); cnt == 0 {
		out.Clear()
	} else {
		out.Grow(cnt)
		for i, el := range els {
			if msg, e := parseMessage(el); e != nil {
				err = e
			} else {
				dst := walk.Walk(out.Index(i))
				err = dec.readMsg(msg, dst)
			}
		}
	}
	return
}

func (dec *Decoder) repeatSlot(slot string, out r.Value, val any) (err error) {
	if els, ok := val.([]any); !ok { // single values can stand in as a slice of one
		err = dec.repeatSlot(slot, out, []any{val})
	} else if cnt := len(els); cnt == 0 {
		out.Clear()
	} else {
		out.Grow(cnt)
		out.SetLen(cnt)
		for i, el := range els {
			if e := dec.slotData(slot, out.Index(i), el); e != nil {
				err = e
				break
			}
		}
	}
	return
}
