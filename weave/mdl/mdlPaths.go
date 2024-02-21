package mdl

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

type paths struct {
	// some ugly caching:
	aspectPath, kindsPath, macroPath, patternPath string // ex. ',4,'
}

func newPaths() *paths {
	return &paths{
		aspectPath:  uncached,
		kindsPath:   uncached,
		macroPath:   uncached,
		patternPath: uncached,
	}
}

// set to something that wont match until its set properly.
const uncached = "$uncached"

func (pen *Pen) getPatternPath() (ret string, err error) {
	return getPath(pen.db, kindsOf.Pattern, &pen.paths.patternPath)
}

func (pen *Pen) getAspectPath() (ret string, err error) {
	return getPath(pen.db, kindsOf.Aspect, &pen.paths.aspectPath)
}

func getPath(db *tables.Cache, kind kindsOf.Kinds, out *string) (ret string, err error) {
	if *out != uncached {
		ret = *out
	} else {
		var path string
		e := db.QueryRow(`
		select (',' || rowid || ',') 
		from mdl_kind where kind = ?1
		limit 1
		`, kind.String()).Scan(&path)
		switch e {
		case nil:
			ret, *out = path, path
		case sql.ErrNoRows:
			err = errutil.New("couldn't determine", kind)
		default:
			err = e
		}
	}
	return
}
