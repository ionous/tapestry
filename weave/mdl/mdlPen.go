package mdl

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/tables"
)

type Pen struct {
	db         *tables.Cache
	paths      *paths
	domain, at string
	warn       Log
}

type Log func(fmt string, parts ...any)

func eatDuplicates(l Log, e error) (err error) {
	if e == nil || !errors.Is(e, Duplicate) {
		err = e
	} else if l != nil {
		l(e.Error())
	}
	return
}
