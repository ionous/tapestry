package rift

import "git.sr.ht/~ionous/tapestry/support/charm"

type CollectionTarget interface {
	// a non-recursive helper to search for the handler of the passed indentation.
	// return the current target if it can handle the indent
	// or return the target's parent if not.
	// see also; PopHistory()
	History(int) CollectionTarget
	// add a new value to the targeted collection,
	// or error if the value isn't acceptable
	WriteValue(any) error
	// read the input and return a new state
	// ( which presumably knows when to WriteValue() and PopHistory() )
	NewRune(rune) charm.State
}

// check whether the passed target can handle the desired level of indentation
// if not -- popup to the parent.
func PopHistory(n CollectionTarget, indent int) (ret charm.State) {
	for {
		// doesn't bother erroring; eventually we'll get to the top
		// and parse will handle things
		if p := n.History(indent); p == n || p == nil {
			ret = n
			break
		} else {
			n = p // try again
		}
	}
	return
}
