package rules

import (
	r "reflect"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/walk"
	"git.sr.ht/~ionous/tapestry/rt"
)

func FilterHasCounter(filter rt.BoolEval) (okay bool) {
	if filter != nil {
		if m, ok := filter.(jsn.Marshalee); !ok {
			panic("unknown type")
		} else {
			okay = searchCounters(m)
		}
	}
	return
}

// fix? could we instead just strstr for countOf
// also might be cool to augment or replace the serialized type
// with our own that has an pre-calced field ( at import, via state parser )
func searchCounters(m jsn.Marshalee) (okay bool) {
	if ok, e := searchForFlow(m, core.CallTrigger_Type); e != nil && e != jsn.Missing {
		panic(e)
	} else {
		okay = ok != nil
	}
	return
}

// return the first flow of the passed type
func searchForFlow(src jsn.Marshalee, typeName string) (ret any, err error) {
	w := walk.Walk(r.ValueOf(src).Elem())
	evts := walk.Callbacks{
		OnFlow: func(w walk.Walker) (err error) {
			if typeName == w.TypeName() {
				ret = w.Value().Interface()
				err = walk.DoneVisiting
				w.Value()
			}
			return
		},
	}
	if e := walk.VisitFlow(w, evts); e != walk.DoneVisiting {
		err = e
	}
	return
}
