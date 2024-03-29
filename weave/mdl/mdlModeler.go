package mdl

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

const Conflict = weaver.Conflict
const Duplicate = weaver.Duplicate
const Missing = weaver.Missing

type FieldInfo = weaver.FieldInfo

/**
 *
 */
func NewModeler(db *sql.DB) (ret *Modeler, err error) {
	return NewModelerWithWarnings(db, nil)
}

func NewModelerWithWarnings(db *sql.DB, warn Log) (ret *Modeler, err error) {
	if warn == nil {
		warn = func(format string, parts ...any) {}
	}
	ret = &Modeler{
		db:    tables.NewCache(db),
		warn:  warn,
		paths: make(paths),
	}
	return
}

// Modeler wraps writing to the model table
// so the implementation can handle verifying dependent names when needed.
type Modeler struct {
	db    *tables.Cache
	paths paths
	warn  Log
}

func (m *Modeler) Pin(domain, at string) *Pen {
	return &Pen{db: m.db, paths: m.paths, domain: domain, at: at, warn: m.warn}
}

// meant for tests which inject their own data outside of weave
func (m *Modeler) PrecachePaths() {
	for _, k := range kindsOf.DefaultKinds {
		m.paths.cachePath(m.db, k)
	}
}
