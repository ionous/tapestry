package rules

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
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
func searchCounters(i jsn.Marshalee) (okay bool) {
	if ok, e := searchForType(i, core.CallTrigger_Type); e != nil && e != jsn.Missing {
		panic(e)
	} else {
		okay = ok
	}
	return
}

func searchForType(src jsn.Marshalee, typeName string) (okay bool, err error) {
	ts := chart.MakeEncoder()
	// fix use panic / recover to early exit?
	var earlyOut error
	err = ts.Marshal(src, &chart.StateMix{
		OnBlock: func(b jsn.Block) error {
			if b.GetType() == typeName {
				okay, earlyOut = true, jsn.Missing
			}
			return earlyOut
		},
		OnKey:   func(_, _ string) error { return earlyOut },
		OnValue: func(_ string, _ interface{}) error { return earlyOut },
	})
	return
}
