package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// implements Assignment for (database) errors
// ( save us querying to produce the same error over and over )
type errorValue struct{ err error }

func (f errorValue) Affinity() affine.Affinity {
	return ""
}

func (f errorValue) GetAssignedValue(run rt.Runtime) (_ g.Value, err error) {
	err = f.err
	return
}
