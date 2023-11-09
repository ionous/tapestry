package rift

import "git.sr.ht/~ionous/tapestry/support/charm"

type Collection interface {
	Document() *Document
	WriteValue(any) error
	Comments() CommentWriter
}

func StartSequence(c *Sequence) charm.State {
	// push to loop back here when we return to this indent.
	doc := c.Document()
	return doc.Push(c.depth,
		charm.Self("sequence", func(loop charm.State, r rune) (ret charm.State) {
			switch r {
			case Hash:
				// fix fix: header comments-- probably should live in collection entries, since its common to all
				panic("not implemented")
			case Dash:
				// every sequence dash lives in the comment block as a vertical tab
				c.comments.WriteRune(VTab)
				ret = loopEntries(c, loop, c.depth)
			}
			return
		}))
}

func StartMapping(c *Mapping) charm.State {
	// push to loop back here when we return to this indent.
	doc := c.Document()
	return doc.Push(c.depth,
		charm.Self("mapping", func(loop charm.State, r rune) (ret charm.State) {
			switch r {
			case Hash:
				// fix fix: header comments
				panic("not implemented")
			default:
				ret = charm.RunStep(r, &c.key, charm.Statement("after key", func(r rune) (ret charm.State) {
					switch r {
					// differentiate brand new keys vs. keys whos values are on a new line.
					// case Newline:
					// 	ret = NextIndent(func() (ret charm.State) {
					// 		if doc.Col <= c.depth {
					// 			ret = loop
					// 		} else {
					// 			ret = loopEntries(c, loop, c.depth)
					// 		}
					// 		return
					// 	})
					default:
						ret = loopEntries(c, loop, c.depth)
					}
					return
				}))
			}
			return
		}))
}

// depth here is the column in which the marker appears
// ( for example: zero if left aligned )
func loopEntries(c Collection, n charm.State, depth int) charm.State {
	//
	// require at least one space after the collection marker
	// ( slightly different than yaml sequences, but i like the consistency of placement )
	//
	// tbd: how will this work with the root document?
	// ( theoretically, it gets one collection entry, so no need for "collection entries" plural. )
	//
	return charm.Step(CollectionEntry(c, depth+2),
		charm.Self("after entry", func(afterEntry charm.State, r rune) (ret charm.State) {
			switch r {
			case Newline:
				// search for the next indent:
				// could be the collection pushed above, or some earlier collection when nested.
				doc := c.Document()
				ret = NextIndent(doc.Pop)
			}
			return
		}))
}
