package rift

import (
	"strings"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// a sequence of array values are specified with:
// a dash, whitespace, the value, trailing whitespace.
// then loops back to itself to handle the next dash.
type Sequence struct {
	doc   *Document
	depth int
	CommentBlock
	values []any
}

// depth is tracked because during value parsing sequences are created after
// determining whether the dash is a minus sign, so the doc position isnt the real position
// alt: create the sequence ahead of time, but would have to handle the push state timing
// ( ex. maybe on first rune? ) and would still have to remember the initial position.
func NewSequence(doc *Document, header string, depth int) *Sequence {
	c := &Sequence{doc: doc, depth: depth}
	if doc.keepComments {
		c.keepComments = true
		c.values = make([]any, 1)
		c.comments.WriteString(header)
	}
	return c
}

// map's empty value is guarded by a completed ke
const emptyValue = errutil.Error("empty value")

// a state that can parse one dash - content pair
// maybe push the returned thingy
func (c *Sequence) NewEntry() charm.State {
	ent := riftEntry{
		doc:          c.doc,
		depth:        c.depth + 2,
		pendingValue: computedValue{emptyValue},
		addsValue: func(val any, comment string) (_ error) {
			if val != emptyValue {
				c.values = append(c.values, val)
			}
			c.comments.WriteString(comment)
			c.comments.WriteRune(Record)
			return
		},
	}
	next := charm.Self("sequence", func(self charm.State, r rune) (ret charm.State) {
		switch r {
		case Hash:
			// this is in between sequence entries
			// potentially, its a header comment for the next element
			// if there is no element, it could be considered a tail
			// of the parent container; it can have nesting.

			ret = charm.RunState(r, HeaderRegion(&ent, c.depth, self))
		case Dash:
			// unlike map, we dont need to hand off the dash itself;
			// only the runes after; also: map's nil value is guarded by a completed key
			// for sequence we have to at least have a dash before we could have a value.
			ent.pendingValue = computedValue{}
			ret = ContentsLoop(&ent)
		}
		return
	})
	return c.doc.PushCallback(ent.depth, next, ent.finalizeEntry)
}

// used by parent collections to read the completed collection
func (c *Sequence) FinalizeValue() (ret any, err error) {
	if c.keepComments {
		comment := strings.TrimRightFunc(c.comments.String(), unicode.IsSpace)
		c.values[0] = comment
	}
	ret, c.values = c.values, nil
	return
}
