package rift

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// after a collection marker,
// read the comment and value ( if either )
//
// [-] <inline spaces> ( <inline comment> | <value> )
//
//	<buffered comment>
//	   <additional buffered lines>
//	<header comment>
//	<value>
func CollectionEntry(c Collection, depth int) charm.State {
	doc := c.Document()
	ent := collectionEntry{Collection: c, indent: depth}
	return doc.PushCallback(depth, &ent, func() error {
		return ent.flush()
	})
}

type collectionEntry struct {
	Collection
	buffer, header strings.Builder
	bufferedLines  int
	indent         int // we could pull this from the doc stack i suppose....
}

func (ent *collectionEntry) flushHeader() (ret string, err error) {
	if ent.bufferedLines > 1 {
		err = errutil.New("ambiguous multiline comment.")
	} else {
		ret = ent.header.String()
		ent.header.Reset()
		ent.flush()
	}
	return
}

func (ent *collectionEntry) flush() (_ error) {
	ent.CommentWriter().WriteString(ent.buffer.String())
	ent.CommentWriter().WriteString(ent.header.String())
	return
}

func (ent *collectionEntry) WriteValue(v any) (_ error) {
	ent.flush()
	return ent.Collection.WriteValue(v)
}

// start reading the collection entry
// ( the padding to the right of a collection marker )
func (ent *collectionEntry) NewRune(r rune) charm.State {
	return charm.Self("padding", func(padding charm.State, r rune) (ret charm.State) {
		switch r {
		case Space:
			ret = padding

		case Hash:
			// these use >= so that content can appear at column zero in documents
			if doc := ent.Document(); doc.Col >= ent.indent {
				ret = ReadComment(ent.CommentWriter(), padding)
			}
		case Newline:
			ret = NextIndent(func() (ret charm.State) {
				if doc := ent.Document(); doc.Col >= ent.indent {
					ret = BufferRegion(ent, doc.Col)
				} else {
					ret = doc.Pop()
				}
				return
			})
		default:
			if doc := ent.Document(); doc.Col >= ent.indent {
				ret = Value(ent, r)
			}
		}
		return
	}).NewRune(r)
}

// we are at the start of a line where comment buffering might occur.
func BufferRegion(ent *collectionEntry, depth int) charm.State {
	return charm.Self("buffering", func(buffering charm.State, r rune) (ret charm.State) {
		switch r {

		// possibly a value at the same depth as the buffering section
		default:
			ret = Value(ent, r)

		// after a completely empty line: move to the header region.
		case Newline:
			ret = HeaderRegion(ent, depth)

			// read comment, and search for next indent
		case Hash:
			ret = ReadComment(&ent.buffer,
				NextIndent(func() (ret charm.State) {
					ent.bufferedLines++
					switch doc := ent.Document(); {
					case doc.Col == depth:
						// at the same indent, stick where we're at.
						ret = buffering
					case ent.bufferedLines == 1 && doc.Col > depth:
						// the ideal multiline comment a single comment followed by indented lines
						ret = IndentedComment(ent, doc.Col)
					default:
						ret = doc.Pop()
					}
					return
				}))
			return
		}
		return
	})
}

// anything at the same indent can be a continuing comment
// anything at a different indent can be the header or value.
func IndentedComment(ent *collectionEntry, depth int) (ret charm.State) {
	out := &ent.buffer
	return charm.Self("indented comment", func(indentedComment charm.State, r rune) (ret charm.State) {
		switch r {
		case Hash:
			ret = ReadComment(out, NextIndent(func() (ret charm.State) {
				// loop if we're at our indent, or start treating future comments as value header comments.
				if doc := ent.Document(); doc.Col == depth {
					ret = indentedComment
				} else {
					ret = HeaderRegion(ent, doc.Col)
				}
				return
			}))

		}
		return
	})
}

// after the buffering section
// collect the header and get ready to pass it to the value.
func HeaderRegion(ent *collectionEntry, depth int) (ret charm.State) {
	out := &ent.header
	return charm.Self("header comment", func(headerComment charm.State, r rune) (ret charm.State) {
		switch r {
		case Hash:
			ret = ReadComment(out, NextIndent(func() (ret charm.State) {
				if doc := ent.Document(); doc.Col == depth {
					ret = headerComment
				} else {
					ret = doc.Pop()
				}
				return
			}))

		default:
			ret = Value(ent, r)
		}
		return
	})
}

// at the start of a rune which might be a value:
func Value(ent *collectionEntry, r rune) (ret charm.State) {
	// dont bother trying to read a value if it wasn't meant to be.
	if r != Newline && r != Space {
		doc := ent.Document()
		pv := &valueState{entry: ent}
		state := doc.PushCallback(doc.Col, pv, func() error {
			return pv.finalizeValue()
		}).NewRune(r)
		// fix ... need to implement this:
		ret = charm.Step(state, charm.Self("trailing comments", func(tail charm.State, r rune) (ret charm.State) {
			switch r {
			case Space:
				ret = tail
			}
			return
		}))
	}
	return
}
