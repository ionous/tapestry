package rift

import (
	"io"
	"strings"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift/maps"
)

type Document struct {
	History
	Cursor
	value any
	CommentBlock
	MakeMap maps.BuilderFactory
}

type Result struct {
	Content any
	// document level comments
	// sub-collection comments are stored with the sequence or mapping
	Comment string
}

func NewDocument(mapMaker maps.BuilderFactory, cmtMaker CommentFactory) *Document {
	return &Document{MakeMap: mapMaker, CommentBlock: cmtMaker()}
}

// has incorrect behavior if called multiple times
func (doc *Document) ReadDoc(src io.RuneReader) (ret Result, err error) {
	if e := doc.ReadLines(src, doc.NewEntry()); e != nil {
		err = e
	} else {
		ret, err = doc.Finalize()
	}
	return
}

// slightly lower level access for reading explicit kinds of values
// calling this multiple times leads to undefined results (fix?)
func (doc *Document) ReadLines(src io.RuneReader, start charm.State) (err error) {
	run := charm.Parallel("parse lines", FilterControlCodes(), UnhandledError(start), &doc.Cursor)
	if e := charm.Read(src, run); e != nil {
		err = e
	} else if e := doc.PopAll(); e != nil {
		err = e
	}
	return
}

// ugly: if preserve comments is true,
// { value, comment, error }
func (doc *Document) Finalize() (ret Result, err error) {
	ret.Content, doc.value = doc.value, nil
	if doc.keepComments {
		ret.Comment = strings.TrimRightFunc(doc.comments.String(), unicode.IsSpace)
	}
	return
}

// create an initial reader state
func (doc *Document) NewEntry() charm.State {
	ent := riftEntry{
		doc:          doc,
		depth:        0,
		pendingValue: computedValue{},
		addsValue: func(val any, comment string) (_ error) {
			doc.value = val // tbd: error if already written?
			doc.comments.WriteString(comment)
			return
		},
	}
	return doc.PushCallback(0, Contents(&ent), ent.finalizeEntry)
}

// pop parser states up to the current indentation level
func (doc *Document) popToIndent() charm.State {
	return doc.History.Pop(doc.Cursor.Col)
}
