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
	if ok, e := searchForFlow(i, core.CallTrigger_Type); e != nil && e != jsn.Missing {
		panic(e)
	} else {
		okay = ok != nil
	}
	return
}

// return the first flow of the passed type
func searchForFlow(src jsn.Marshalee, typeName string) (ret any, err error) {
	ts := chart.MakeEncoder()
	var earlyOut error // tbd: use panic / recover to early exit more quickly?
	err = ts.Marshal(src, &chart.StateMix{
		OnKey:   func(_, _ string) error { return earlyOut },
		OnValue: func(_ string, _ interface{}) error { return earlyOut },
		OnBlock: func(b jsn.Block) error {
			if earlyOut == nil && b.GetType() == typeName {
				if flow, ok := b.(jsn.FlowBlock); ok {
					if op := flow.GetFlow(); op != nil {
						ret, earlyOut = op, jsn.Missing
					}
				}
			}
			return earlyOut
		},
	})
	return
}
