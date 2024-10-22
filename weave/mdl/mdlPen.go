package mdl

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/tables"
)

type Pen struct {
	db     *tables.Cache
	paths  paths
	domain string
	pos    compact.Source
	warn   Log
}

type Log func(fmt string, parts ...any)

// if the error is from a duplicate definition, return nil as if there was no error.
// log duplicate errors if the logger is valid.
func eatDuplicates(l Log, e error) (err error) {
	if e == nil || !errors.Is(e, ErrDuplicate) {
		err = e
	} else if l != nil {
		l(e.Error())
	}
	return
}
