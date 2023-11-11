package rift

import (
	"io"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift/maps"
)

type Document struct {
	History
	Cursor
	Value any
	CommentBlock
	MakeMap maps.BuilderFactory
}

func NewDocument(mapMaker maps.BuilderFactory, cmtMaker CommentFactory) *Document {
	return &Document{MakeMap: mapMaker, CommentBlock: cmtMaker()}
}

// cant be called multiple times (fix?)
func (doc *Document) ReadDoc(src io.RuneReader) error {
	return doc.ReadLines(src, doc.NewEntry())
}

// slightly lower level access for reading explicit kinds of values
// cant be called multiple times (fix?)
func (doc *Document) ReadLines(src io.RuneReader, start charm.State) (err error) {
	run := charm.Parallel("parse lines", FilterControlCodes(), UnhandledError(start), &doc.Cursor)
	if e := charm.Read(src, run); e != nil {
		err = e
	} else if e := doc.PopAll(); e != nil {
		err = e
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
			doc.Value = val // tbd: error if already written?
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
