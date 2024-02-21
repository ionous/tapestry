package mdl

import (
	"database/sql"
	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/tapestry/tables"
)

// when the definition would contradict existing information:
// the returned error wraps this tag. errors.Is can be used to detect it.
const Conflict = errutil.Error("Conflict")

// when the definition would repeat existing information:
// the returned error wraps this tag. errors.Is can be used to detect it.
const Duplicate = errutil.NoPanicError("Duplicate")

// when the definition can't find some required information:
// the returned error wraps this tag. errors.Is can be used to detect it.
const Missing = errutil.NoPanicError("Missing")

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
		paths: newPaths(),
	}
	return
}

// Modeler wraps writing to the model table
// so the implementation can handle verifying dependent names when needed.
type Modeler struct {
	db    *tables.Cache
	paths *paths
	warn  Log
}

func (m *Modeler) Pin(domain, at string) *Pen {
	return &Pen{db: m.db, paths: m.paths, domain: domain, at: at, warn: m.warn}
}
