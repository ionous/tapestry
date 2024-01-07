package encode

import "git.sr.ht/~ionous/tapestry/lang/walk"

func (enc *Encoder) encodeFlows(it walk.Walker) (ret []any, err error) {
	if cnt := it.Len(); cnt > 0 {
		ret = make([]any, cnt)
		for i := 0; i < cnt && it.Next(); i++ {
			if res, e := enc.writeFlow(it); e != nil {
				ret, err = nil, e
				break
			} else {
				ret[i] = res
			}
		}
	}
	return
}

func (enc *Encoder) encodeSlots(it walk.Walker) (ret []any, err error) {
	if cnt := it.Len(); cnt > 0 {
		ret = make([]any, cnt)
		for i := 0; it.Next(); i++ { // the indexed element
			// walk, the container of the slot
			// next, its possibly empty contents
			if slot := it.Walk(); slot.Next() {
				if res, e := enc.writeFlow(slot.Walk()); e != nil {
					ret, err = nil, e
					break
				} else {
					ret[i] = res
				}
			}
		}
	}
	return
}
