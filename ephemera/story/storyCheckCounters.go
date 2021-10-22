package story

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
)

func searchForCounters(i jsn.Marshalee) (okay bool) {
	if ok, e := searchForType(i, core.CallTrigger_Type); e != nil {
		panic(e)
	} else {
		okay = ok
	}
	return
}

func searchForType(src jsn.Marshalee, typeName string) (okay bool, err error) {
	ts := chart.MakeEncoder(nil)
	// fix use panic / recover to early exit?
	err = ts.Marshal(src, &chart.StateMix{
		OnWarn: ts.Error,
		OnBlock: func(b jsn.BlockType) bool {
			if b.GetType() == typeName {
				okay = true
			}
			return !okay
		},
		OnKey: func(_, _ string) bool {
			return !okay
		},
		OnValue: func(_ string, _ interface{}) {},
		OnEnd:   func() {},
	})
	return
}
