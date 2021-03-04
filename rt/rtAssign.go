package rt

import (
	"git.sr.ht/~ionous/iffy/affine"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

//  Assignment - limits variable and parameter assignment to particular contexts.
type Assignment interface {
	// fix? used by assembly decodeProg and story importer importAgs
	// to verify the affinity of local initializers and arguments.
	Affinity() affine.Affinity
	// write the results of evaluating this into that.
	GetAssignedValue(Runtime) (g.Value, error)
}

type Arg struct {
	Name string
	From Assignment
}
