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
	comments strings.Builder
	values   []any
}

// maybe doc is a factory even?
func NewSequence(doc *Document, indent int, writeBack func(vs []any) error) charm.State {
	// create the values list, always saving space for a comment block
	seq := &Sequence{doc: doc, values: []any{nil}}
	return doc.PushIndent(indent, seq, func() error {
		// write the comment block
		comment := strings.TrimRightFunc(seq.comments.String(), unicode.IsSpace)
		seq.values[0] = comment
		return writeBack(seq.values)
	})
}

// called for every sequence element at a given nesting level
func (n *Sequence) NewRune(first rune) (ret charm.State) {
	switch first {
	case Comment:
		// in theory this can only happen after having at least one sequence dash.
		ret = parseCollectionComments(n.doc, &n.comments)

	case SequenceDash:
		// every sequence dash lives in the comment block as a vertical tab
		n.comments.WriteRune(VTab)

		// cheating a bit here: if the next rune is whitespace or an eof
		// then there will be no pop() to write a value, so write it here
		// alt: push back a pending state into the history, or track a bool
		n.values = append(n.values, nil)
		ret = parseCollection(n.doc, func(val any) (_ error) {
			// and rewrite the value here:
			n.values[len(n.values)-1] = val
			return
		})
	}
	return
}
