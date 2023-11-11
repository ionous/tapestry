package rift

import (
	"strings"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/charm"
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
		if len(header) > 0 {
			c.comments.WriteString(header)
			if len(header) > 0 {
				c.comments.WriteString(header)
			}
		}
	}
	return c
}

// a state that can parse one dash - content pair
// maybe push the returned thingy
func (c *Sequence) NewEntry() charm.State {
	ent := riftEntry{
		doc:          c.doc,
		depth:        c.depth + 2,
		pendingValue: computedValue{},
		addsValue: func(val any, comment string) (_ error) {
			c.values = append(c.values, val)
			c.comments.WriteString(comment)
			return
		},
	}
	next := charm.Statement("sequence", func(r rune) (ret charm.State) {
		switch r {
		case Hash:
			// fix fix: header comments-- probably should live in collection entries, since its common to all
			panic("not implemented")

		case Dash:
			// every sequence dash lives in the comment block as a vertical tab
			c.comments.WriteRune(Carriage)
			// unlike map, we dont need to hand off the dash itself;
			// only the runes after.
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
