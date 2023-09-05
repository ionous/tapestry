package rules

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
)

// words that can precede a rule name
// ex. "instead of taking"
type eventPrefix int

//go:generate stringer -type=eventPrefix -linecomment
const (
	before  eventPrefix = iota
	instead             // instead of
	after
	report
	//
	numPrefixes = iota
)

// with the theory that sqlite sorts asc by default
// ( spaces numbers out so some theoretical extra 'first' or 'last' prefix could add or subtract to the numbers )
func (p eventPrefix) rank() (ret int) {
	return ranks[p]
}

// each rule gets a sort order group, smaller numbers run earlier than larger numbers.
// within each group the last ( most recently defined ) rule wins.
var ranks = []int{-20, -10, 10, 20, 0}

// return name sans any prefix, and any prefix the name had.
func findPrefix(name string) (short string, prefix eventPrefix) {
	short, prefix = name, numPrefixes // provisional
	for i := 0; i < numPrefixes; i++ {
		p := eventPrefix(i).String()
		if strings.HasPrefix(name, p+" ") {
			short = name[len(p)+1:]
			prefix = eventPrefix(i)
			break
		}
	}
	return
}

// note: stopping before the action happens is considered a cancel.
func (p eventPrefix) stopJump() (stop bool, jump rt.Jump) {
	switch p {
	case before:
		// before falls through
		stop, jump = true, rt.JumpLater
	case instead:
		// instead stops immediately
		stop, jump = true, rt.JumpNow
	default:
		// a normal rule doesnt stop, it moves on to the next phase
		stop, jump = false, rt.JumpNow
	case after:
		// after falls through
		stop, jump = false, rt.JumpLater
	case report:
		// report stops immediately
		stop, jump = true, rt.JumpNow
	}
	return
}

func (p eventPrefix) eventName(long, short string) (ret string) {
	switch p {
	case instead:
		ret = event.BeforePhase.PatternName(short)
	case report:
		ret = event.AfterPhase.PatternName(short)
	default:
		ret = long
	}
	return
}
