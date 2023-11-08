package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
)

// read everything until the end of the line as a comment.
func ReadComment(out RuneWriter, eol charm.State) charm.State {
	return charm.Self("read comment", func(self charm.State, r rune) (ret charm.State) {
		if r == Newline {
			ret = eol
		} else {
			out.WriteRune(r)
			ret = self
		}
		return
	})
}
