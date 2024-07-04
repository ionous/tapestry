package mdl

import (
	"errors"
	"path"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/tables"
)

type Pen struct {
	db     *tables.Cache
	paths  paths
	domain string
	pos    Source
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

func MakeSource(t typeinfo.Markup) Source {
	var pos Source
	if t != nil {
		m := t.GetMarkup(false)
		if at, ok := m[compact.Position].([]int); ok {
			pos.Line = at[1]
		}
		if at, ok := m[compact.Source].(string); ok {
			file := path.Base(at) // extract the file from shared/something.tell
			pos.File = file
			if full, part := len(at), len(file); full > part {
				pos.Path = at[:full-(part+1)] // skip trailing slash before the filename
			}
		}
	}
	return pos
}
