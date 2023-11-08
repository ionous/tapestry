package rift

import "git.sr.ht/~ionous/tapestry/support/charm"

// find the next indent and the state to handle it once we are there.
// ( and send the actual next rune to the target's returned state )
func NextIndent(onIndent func() charm.State) (ret charm.State) {
	return charm.Self("next indent", func(findNext charm.State, r rune) (ret charm.State) {
		if r == Space || r == Newline {
			ret = findNext
		} else if next := onIndent(); isDone(next) {
			ret = next
		} else {
			ret = next.NewRune(r)
		}
		return
	})
}

// return true if the passed state is unhandled or in error
func isDone(c charm.State) (okay bool) {
	switch c.(type) {
	case nil, charm.Terminal:
		okay = true
	}
	return
}
