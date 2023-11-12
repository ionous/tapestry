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

// nested comments are fixed at the passed depth
// starts on something other than whitespace
func NestedComment(doc *Document, out *CommentBuffer) charm.State {
	depth := doc.Col
	return charm.Self("nested comment", func(self charm.State, r rune) (ret charm.State) {
		switch r {
		case Hash:
			out.WriteRune(Nestline)
			ret = ReadComment(out, self)
		case Newline:
			ret = MaintainIndent(doc, self, depth)
		}
		return
	})
}
