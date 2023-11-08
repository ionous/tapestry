package rift

import "git.sr.ht/~ionous/tapestry/support/charm"

type RuneWriter interface {
	WriteRune(r rune) (int, error)
	WriteString(s string) (int, error)
}

type Collection interface {
	Document() *Document
	WriteValue(any) error
	CommentWriter() RuneWriter
}

func CollectionEntries(c Collection, n charm.State, depth int) charm.State {
	// require at least one space after the collection marker
	// ( slightly different than yaml sequences, but i like the consistency of placement )
	return charm.Step(CollectionEntry(c, depth+1),
		charm.Self("after entry", func(afterEntry charm.State, r rune) (ret charm.State) {
			switch r {
			case Newline:
				// if no one handled the end of line, then it probably indicates the end of the entry.
				// search for the next ident.
				ret = NextIndent(c.Document().Pop)
			}
			return
		}))
}
