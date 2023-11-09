package rift

import "git.sr.ht/~ionous/tapestry/support/charm"

type Collection interface {
	Document() *Document
	WriteValue(any) error
	Comments() CommentWriter
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
