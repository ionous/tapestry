package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

type Document struct {
	Value any
}

// when we're done, we're done.
func (d *Document) History(indent int) (null CollectionTarget) {
	return
}

// parse a single value
func (d *Document) NewRune(r rune) charm.State {
	value := ValueParser{d, 0}
	return value.NewRune(r)
}

// receives the value from the value parser
func (d *Document) WriteValue(val any) (err error) {
	if d.Value != nil {
		err = errutil.New("document already has a value")
	} else {
		d.Value = val
	}
	return
}
