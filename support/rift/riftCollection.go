package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// all collection value handling is the same:
// look for whitespace, check the indentation level to see if there was a null value,
// read the value ( including any sub-collections ) and find the next newline,
// check the indentation to see if we're sticking with the same collection, or popping to a previous one.
func parseCollection(doc *Document, onValue func(any) error) (ret charm.State) {
	startingIndent := doc.CurrentIndent() // padding is the space between the dash or colon and the value

	return CommentSpaces("padding", startingIndent, func(padding int) (ret charm.State) {
		switch {
		// if the indent is less or equal, than we're a new line.
		// this returns us to our sequence, mapping, or document state.
		// ( the value was already written nil after the dash )
		case padding <= startingIndent:
			ret = doc.PopIndent(padding)
		// an increased indentation means a value:
		default:
			// read the value ( could be another collection )
			ret = charm.Step(NewValue(doc, padding, onValue), charm.Statement("after value", func(r rune) charm.State {

				// AFTER THE VALUE THERE CAN BE COMMENTS
				// this is where we'd write the \t

				// after value we require a newline:
				return charm.RunState(r, CommentLines("tail", padding, func(tail int) (ret charm.State) {
					switch {
					case tail <= startingIndent:
						ret = doc.PopIndent(tail)
					default:
						// the only valid increases are collections,
						// and after a collection, we should be at a less or equal indent.
						e := errutil.New("invalid indentation")
						ret = charm.Error(e)
					}
					return
				}))
			}))
		}
		return
	})
}
