package decode

import (
	"git.sr.ht/~ionous/tapestry/lang/inspect"
)

func (dec *Decoder) repeatFlow(out inspect.Iter, val any) (err error) {
	if els, ok := val.([]any); !ok { // single values can stand in as a slice of one
		err = dec.repeatFlow(out, []any{val})
	} else {
		out.Resize(len(els))
		for i := 0; out.Next(); i++ {
			if msg, e := ParseMessage(els[i]); e != nil {
				err = e
			} else {
				err = dec.readMsg(msg, out.Walk())
			}
		}
	}
	return
}

func (dec *Decoder) repeatSlot(out inspect.Iter, val any, slotName string) (err error) {
	if els, ok := val.([]any); !ok { // single values can stand in as a slice of one
		err = dec.repeatSlot(out, []any{val}, slotName)
	} else {
		out.Resize(len(els))
		for i := 0; out.Next(); i++ {
			if e := dec.decodeSlot(out.Walk(), els[i], slotName); e != nil {
				err = e
				break
			}
		}
	}
	return
}
