package rules

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/rt"
)

// checks for whether a rule has any counters.
// counters need to be updated whenever a pattern's rules are evaluated
// even if a more important rule has already been matched.
type updateTracker int

// true if any of Check(s) have returned true
func (b updateTracker) HasUpdate() bool {
	return b > 0
}

func (b *updateTracker) CheckArgs(args []assign.Arg) (okay bool) {
	for _, arg := range args {
		slot := rtti.Assignment_Slot{Value: arg.Value}
		if searchCounters(&slot) {
			okay = true
			*b++
			break
		}
	}
	return
}

func (b *updateTracker) CheckFilter(filter rt.BoolEval) (okay bool) {
	if FilterHasCounter(filter) {
		okay = true
		*b++
	}
	return
}
