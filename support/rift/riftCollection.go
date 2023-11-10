package rift

import "git.sr.ht/~ionous/tapestry/support/charm"

// abstraction for reading sets of values
// implemented by documents, mappings, and sequences.
type Collection interface {
	// access to the owner of the collection
	Document() *Document
	// access to writing comments into the collection
	CommentWriter() CommentWriter
	writeValue(any) error
}

func StartSequence(c *Sequence) charm.State {
	doc := c.Document()
	return doc.Push(c.depth, charm.Statement("start sequence", func(r rune) charm.State {
		return c.NewEntry().NewRune(r)
	}))
}

func StartMapping(c *Mapping) charm.State {
	doc := c.Document()
	return doc.Push(c.depth, charm.Statement("start mapping", func(r rune) charm.State {
		return c.NewEntry().NewRune(r)
	}))
}
