package encode

import "git.sr.ht/~ionous/tapestry/lang/walk"

func (enc *Encoder) encodeFlows(it walk.Walker) (ret []any, err error) {
	if cnt := it.ContainerLen(); cnt > 0 {
		ret = make([]any, cnt)
		for i := 0; i < cnt && it.Next(); i++ {
			var sub FlowBuilder
			if e := enc.writeFlow(&sub, it); e != nil {
				ret, err = nil, e
				break
			} else {
				ret[i] = sub.FinalizeMap()
			}
		}
	}
	return
}

func (enc *Encoder) encodeSlots(it walk.Walker) (ret []any, err error) {
	if cnt := it.ContainerLen(); cnt > 0 {
		ret = make([]any, cnt)
		for i := 0; it.Next(); i++ { // the indexed element
			// walk, the container of the slot
			// next, its possibly empty contents
			if slot := it.Walk(); slot.Next() {
				var sub FlowBuilder
				if e := enc.writeFlow(&sub, slot.Walk()); e != nil {
					ret, err = nil, e
					break
				} else {
					ret[i] = sub.FinalizeMap()
				}
			}
		}
	}
	return
}
