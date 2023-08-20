package rules

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/rt"
)

// compiles the results of multiple Check(s)
type updateTracker int

// true if any of Check(s) have returned true
func (b updateTracker) HasUpdate() bool {
	return b > 0
}
func (b *updateTracker) CheckFilter(filter rt.BoolEval) (okay bool) {
	if FilterHasCounter(filter) {
		okay = true
		*b++
	}
	return
}

func (b *updateTracker) CheckArgs(args []assign.Arg) (okay bool) {
	for _, arg := range args {
		if m, ok := arg.Value.(jsn.Marshalee); !ok {
			panic("unknown type")
		} else if searchCounters(m) {
			okay = true
			*b++
			break
		}
	}
	return
}
