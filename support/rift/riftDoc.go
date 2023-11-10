package rift

import (
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift/maps"
	"github.com/ionous/errutil"
)

type Document struct {
	History
	Cursor
	Value any
	CommentBlock
	MakeMap maps.BuilderFactory
}

func (doc *Document) ParseLines(str string, start charm.State) (err error) {
	run := charm.Parallel("parse lines", FilterControlCodes(), UnhandledError(start), &doc.Cursor)
	if e := charm.Parse(str, run); e != nil {
		err = e
	} else if e := doc.PopAll(); e != nil {
		err = e
	}
	return
}

func (doc *Document) Pop() charm.State {
	return doc.History.Pop(doc.Cursor.Col)
}

func (doc *Document) Document() *Document {
	return doc
}

// fix: return error if already written
func (doc *Document) WriteValue(val any) (_ error) {
	doc.Value = val
	return
}

// turns any unhandled states returned by the watched state into errors
func UnhandledError(watch charm.State) charm.State {
	return charm.Self("unhandled error", func(self charm.State, r rune) (ret charm.State) {
		if next := watch.NewRune(r); next == nil {
			ret = charm.Error(errutil.Fmt("unexpected character %q(%d) during %s", r, r, charm.StateName(watch)))
		} else {
			ret, watch = self, next // keep checking until watch returns nil
		}
		return
	})
}

// except for newline, control codes are considered invalid.
func FilterControlCodes() charm.State {
	return charm.Self("filter control codes", func(next charm.State, r rune) charm.State {
		if r != Newline && unicode.IsControl(r) {
			e := errutil.New("invalid character", int(r))
			next = charm.Error(e)
		}
		return next
	})
}
