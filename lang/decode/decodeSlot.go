package decode

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

func (dec *Decoder) decodeSlot(w inspect.Iter, slot *typeinfo.Slot, data any) (err error) {
	if ptr, e := dec.innerDecode(w, slot, data); e != nil {
		err = e
	} else if !w.SetSlot(ptr) {
		err = fmt.Errorf("couldnt assign value of %T to %s", ptr, w.TypeInfo().TypeName())
	}
	return
}

func (dec *Decoder) innerDecode(w inspect.Iter, slot *typeinfo.Slot, data any) (ret typeinfo.Inspector, err error) {
	if v, e := dec.customDecode(w, slot, data); !compact.IsUnhandled(e) {
		ret, err = v, e
	} else {
		if msg, e := ParseMessage(data); e != nil {
			err = e
		} else if ptr, ok := dec.signatures.Create(slot.Name, msg.Key); ok {
			// create the flow, and decode the message into it.
			if e := dec.readMsg(msg, inspect.Walk(ptr)); e != nil {
				err = e
			} else {
				ret = ptr
			}
		} else if dec.patterns != nil {
			// otherwise, try as a pattern
			ret, err = dec.patterns(dec, slot, msg)
		} else {
			// otherwise, an error
			err = fmt.Errorf("unknown signature slot=%q flow=%q", slot, msg.Key)
		}
	}
	return
}

func (dec *Decoder) customDecode(w inspect.Iter, slot *typeinfo.Slot, arg any) (ret typeinfo.Inspector, err error) {
	if c := dec.customDecoder; c == nil {
		err = compact.Unhandled("custom decoder")
	} else {
		ret, err = c(dec, slot, arg)
	}
	return
}
