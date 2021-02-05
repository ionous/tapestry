package term

import (
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// Preparer - terms are stored as part of a pattern.
// they can add their names and default values to the list of expected parameters or predetermined locals
// stored lists of preparers may be replaced by kinds at some point
type Preparer interface {
	Prepare(rt.Runtime) (g.Value, error)
	String() string
}
