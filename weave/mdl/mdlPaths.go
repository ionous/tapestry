package mdl

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

// some ugly caching:
type paths map[kindsOf.Kinds]pathEntry

type pathEntry struct {
	path  string
	error error
}

// set to something that wont match until its set properly.
const uncached = "$uncached"

func (pen *Pen) getPatternPath() (ret string, err error) {
	return pen.paths.cachePath(pen.db, kindsOf.Pattern)
}

func (pen *Pen) getAspectPath() (ret string, err error) {
	return pen.paths.cachePath(pen.db, kindsOf.Aspect)
}

// doesnt cache, doesnt check for errors
func (pen *Pen) getPath(k kindsOf.Kinds) (ret string) {
	if at := pen.paths[k]; len(at.path) > 0 {
		ret = at.path
	} else {
		ret = uncached
	}
	return
}

func (p *paths) cachePath(db *tables.Cache, k kindsOf.Kinds) (ret string, err error) {
	if at, ok := (*p)[k]; ok {
		ret, err = at.path, at.error
	} else {
		ret, err = queryPath(db, k)
		(*p)[k] = pathEntry{ret, err}
	}
	return
}

func queryPath(db *tables.Cache, kind kindsOf.Kinds) (ret string, err error) {
	switch e := db.QueryRow(`
		select (',' || rowid || ',') 
		from mdl_kind where kind = ?1
		limit 1
		`, kind.String()).Scan(&ret); e {
	case sql.ErrNoRows:
		err = errutil.New("couldn't determine", kind)
	default:
		err = e // nil or ortherwise
	}
	return
}
