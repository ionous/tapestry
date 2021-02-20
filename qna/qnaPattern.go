package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

// fix: for GetEvalByName. ideally this goes away
type patternValue struct {
	store interface{}
}

func (f patternValue) Affinity() affine.Affinity {
	return "" // not needed currently
}

func (f patternValue) GetAssignedValue(run rt.Runtime) (_ g.Value, err error) {
	err = errutil.New("pattern expected use of GetEvalByName")
	return
}
