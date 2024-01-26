package rules

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
)

func FilterHasCounter(filter rt.BoolEval) (okay bool) {
	if filter != nil {
		slot := rtti.BoolEval_Slot{Value: filter}
		okay = searchCounters(&slot)
	}
	return
}

// fix? could we instead just strstr for countOf
// also might be cool to augment or replace the serialized type
// with our own that has an pre-calced field ( at import, via state parser )
func searchCounters(m typeinfo.Inspector) (okay bool) {
	if ok, e := searchForFlow(m, &core.Zt_CallTrigger); e != nil {
		panic(e)
	} else {
		okay = ok != nil
	}
	return
}

// return the first flow of the passed type
func searchForFlow(src typeinfo.Inspector, find typeinfo.T) (ret any, err error) {
	evts := inspect.Callbacks{
		OnFlow: func(w inspect.Iter) (err error) {
			if find == w.TypeInfo() {
				ret = w.GoValue()
				err = inspect.DoneVisiting
			}
			return
		},
	}
	if e := inspect.Visit(src, evts); e != inspect.DoneVisiting {
		err = e
	}
	return
}
