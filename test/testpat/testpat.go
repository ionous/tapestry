package testpat

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/test/testutil"
	"github.com/ionous/errutil"
)

type Runtime struct {
	pattern.Map
	testutil.Runtime
}

func (run *Runtime) Call(name string, aff affine.Affinity, args []rt.Arg) (ret g.Value, err error) {
	if pat, ok := run.Map[name]; !ok {
		err = errutil.New("unknown pattern", name)
	} else {
		ret, err = pat.Run(run, aff, args)
	}
	return
}
