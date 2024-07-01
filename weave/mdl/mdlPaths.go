package mdl

import (
	"database/sql"

	"fmt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
)

// some ugly caching:
type paths map[kindsOf.Kinds]pathEntry

type pathEntry struct {
	path  string
	error error
}

// set to something that wont match until its set properly.
const uncached = "$uncached"

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
		err = fmt.Errorf("couldn't determine kind %q", kind)
	default:
		err = e // nil or ortherwise
	}
	return
}
