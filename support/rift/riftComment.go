package rift

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// reads all of the comments of a collection which start at the same nesting level
func parseCollectionComments(doc *Document, out *strings.Builder) (ret charm.State) {
	starting := doc.Cursor // our starting position
	out.WriteRune(Comment) // known to start with a comment marker
	return charm.Self("trailing comment", func(trailingComment charm.State, r rune) (ret charm.State) {
		switch r {
		// add everything to the comment
		default:
			out.WriteRune(r)
			ret = trailingComment

		// and after newlines, scan for the first rune that isnt whitespace
		// so that we can properly handle dedents
		case Newline:
			ret = charm.Self("exit comment", func(self charm.State, r rune) (ret charm.State) {
				switch r {
				// eat whitespace
				case Space, Newline:
					ret = self

				// any other rune:
				default:
					switch {
					// another comment? we can read it ( by looping ) even on an increased indent
					// ( tbd: add the extra spaces to the comment? )
					case r == Comment && (doc.Cursor.Indent >= starting.Indent):
						out.WriteRune(Comment)
						ret = trailingComment

					//  less or equal? return to the appropriate collection
					// ( and that's true even for a comment, because it belongs to that other collection )
					case doc.Cursor.Indent <= starting.Indent:
						parentCollection := doc.PopIndent(doc.Cursor.Indent)
						ret = parentCollection.NewRune(r) // dont forget to parse the rune

					default:
						// any other rune at an increased indentation is invalid.
						e := errutil.New("invalid indentation")
						ret = charm.Error(e)
					}
				}
				return
			})

		}
		return
	})
}
