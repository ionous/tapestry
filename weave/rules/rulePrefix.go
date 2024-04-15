package rules

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
)

// words preceding a rule name control the grouping, sorting, and interaction of similar rules.
// ex. "instead of taking" sorts after "before taking" within the "before" grouping.
// "when" is the default
type Prefix int

//go:generate stringer -type=Prefix -linecomment
const (
	Before  Prefix = iota // before
	Instead               // instead of
	When                  // when
	After                 // after
	Report                // report
	//
	NumPrefixes = iota
)

// with the theory that sqlite sorts asc by default
// ( spaces numbers out so some theoretical extra 'first' or 'last' prefix could add or subtract to the numbers )
func (p Prefix) rank() (ret int) {
	return ranks[p]
}

// each rule gets a sort order group, smaller numbers run earlier than larger numbers.
// within each group the last ( most recently defined ) rule wins.
var ranks = []int{-20, -10, 0, 10, 20, 0}

// controls whether a rule blocks the processing of other rules.
// note: stopping before the action happens is considered a cancel.
func (p Prefix) stopJump() (stop bool, jump rt.Jump) {
	switch p {
	case Before:
		// before falls through
		stop, jump = true, rt.JumpLater
	case Instead:
		// instead stops immediately
		stop, jump = true, rt.JumpNow
	case When:
		// a normal rule jumps but it doesnt stop processing entirely;
		// it moves on to the next phase
		stop, jump = false, rt.JumpNow
	case After:
		// after falls through
		stop, jump = false, rt.JumpLater
	case Report:
		// report stops immediately
		stop, jump = true, rt.JumpNow
	}
	return
}

// return name sans any prefix, and any prefix the name had.
func findPrefix(name string) (short string, prefix Prefix) {
	short, prefix = name, When // provisional
	for i := 0; i < NumPrefixes; i++ {
		p := Prefix(i).String()
		if strings.HasPrefix(name, p+" ") {
			short = name[len(p)+1:]
			prefix = Prefix(i)
			break
		}
	}
	return
}
