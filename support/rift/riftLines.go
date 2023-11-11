package rift

import "git.sr.ht/~ionous/tapestry/support/charm"

// read everything until the end of the line as a comment.
// send the newline to the passed state.
func ReadComment(out CommentWriter, eol charm.State) charm.State {
	out.WriteRune(Hash)
	return charm.Self("read comment", func(self charm.State, r rune) (ret charm.State) {
		if r == Newline {
			ret = eol.NewRune(r)
		} else {
			out.WriteRune(r)
			ret = self
		}
		return
	})
}

// find the next indent and the state to handle it once we are there.
// ( and send the actual next rune to the target's returned state )
func NextIndent(onIndent func() charm.State) (ret charm.State) {
	return charm.Self("next indent", func(nextIndent charm.State, r rune) (ret charm.State) {
		if r == Space || r == Newline {
			ret = nextIndent
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
