package rift

import (
	"strings"
	"unicode"
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
func NewSequence(parent Collection, header string, depth int) *Sequence {
	doc := parent.Document()
	c := &Sequence{doc: doc, depth: depth}
	if doc.keepComments {
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

func (c *Sequence) Document() *Document {
	return c.doc
}

func (c *Sequence) WriteValue(val any) (_ error) {
	c.values = append(c.values, val)
	return
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
