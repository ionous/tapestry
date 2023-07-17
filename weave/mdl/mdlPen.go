package mdl

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/tables"
)

/**
 *
 */
func NewModeler(db *sql.DB) (ret *Modeler, err error) {
	ret = &Modeler{
		db: tables.NewCache(db),
		paths: &paths{
			aspectPath:  uncached,
			kindPath:    uncached,
			macroPath:   uncached,
			patternPath: uncached,
		},
	}
	return
}

// Modeler wraps writing to the model table
// so the implementation can handle verifying dependent names when needed.
type Modeler struct {
	db    *tables.Cache
	paths *paths
}

type paths struct {
	// some ugly caching:
	aspectPath, kindPath, macroPath, patternPath string // ex. ',4,'
}

// set to something that wont match until its set properly.
const uncached = "$uncached"

func (m *Modeler) Pin(domain, at string) *Pen {
	return &Pen{db: m.db, paths: m.paths, domain: domain, at: at}
}
