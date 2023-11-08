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
	doc      *Document
	depth    int
	comments strings.Builder
	values   []any
}

// depth is tracked because during value parsing sequences are created after
// determining whether the dash is a minus sign, so the doc position isnt the real position
// alt: create the sequence ahead of time, but would have to handle the push state timing
// ( ex. maybe on first rune? ) and would still have to remember the initial position.
func NewSequence(parent Collection, header string, depth int) *Sequence {
	doc := parent.Document()
	// create the values list, always saving space for a comment block
	c := &Sequence{doc: doc, depth: depth}
	keepComments := doc.KeepComments
	if keepComments {
		c.values = make([]any, 1)
	}
	if len(header) > 0 {
		c.comments.WriteString(header)
	}
	doc.PushCallback(depth, c, func() error {
		// write the comment block
		if keepComments {
			comment := strings.TrimRightFunc(c.comments.String(), unicode.IsSpace)
			c.values[0] = comment
		}
		return parent.WriteValue(c.values)
	})
	return c
}

func (c *Sequence) Document() *Document {
	return c.doc
}

// fix: return error if already written
func (c *Sequence) WriteValue(val any) (_ error) {
	c.values[len(c.values)-1] = val
	return
}

func (c *Sequence) CommentWriter() RuneWriter {
	return &c.comments
}

// called for every sequence element at a given nesting level
func (c *Sequence) NewRune(r rune) (ret charm.State) {
	if r == Hash {
		// fix fix: header comments-- probably should live in collection entries, since its common to all
		panic("not implemented")
	} else if r == Dash {
		// every sequence dash lives in the comment block as a vertical tab
		c.comments.WriteRune(VTab)
		// cheating a bit here: if the next rune is whitespace or an eof
		// then there will be no pop() to write a value, so write it here
		// alt: push back a pending state into the history, or track a bool
		c.values = append(c.values, nil)
		ret = CollectionEntries(c, c, c.depth+1)
	}
	return
}
