package rt

import (
	"git.sr.ht/~ionous/iffy/affine"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

//  Assignment - limits variable and parameter assignment to particular contexts.
type Assignment interface {
	// fix: needed by importArgs right now for ... reasons...
	Affinity() affine.Affinity
	// write the results of evaluating this into that.
	GetAssignedValue(Runtime) (g.Value, error)
}

type Arg struct {
	Name string
	From Assignment
}
