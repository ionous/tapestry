package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// all collection value handling is the same:
// look for whitespace, check the indentation level to see if there was a null value,
// read the value ( including sub-collections ) and find the next newline,
// check the indentation to see if we're sticking with the same collection, or popping to a previous one.
func parseCollection(h *Document, onValue func(any) error) (ret charm.State) {
	startingIndent := h.CurrentIndent() // padding is the space between the dash or colon and the value

	return RequireSpaces("padding", startingIndent, func(padding int) (ret charm.State) {
		switch {
		// if the indent is less or equal,
		// than we're a new line and the value was null
		case padding <= startingIndent:
			ret = h.PopIndent(padding)
		// an increased indentation means a value:
		default:
			// read the value ( could be another collection )
			ret = charm.Step(NewValue(h, padding, onValue), charm.Statement("after value", func(r rune) charm.State {

				// after value we require a newline:
				return charm.RunState(r, RequireLines("tail", padding, func(tail int) (ret charm.State) {
					switch {
					case tail <= startingIndent:
						ret = h.PopIndent(tail)
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
