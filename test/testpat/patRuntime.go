package testpat

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/pattern"
	"git.sr.ht/~ionous/iffy/test/testutil"
	"github.com/ionous/errutil"
)

type Runtime struct {
	Map
	testutil.Runtime
}

func (run *Runtime) Call(name string, aff affine.Affinity, args []rt.Arg) (ret g.Value, err error) {
	if pat, ok := run.Map[name]; !ok {
		err = errutil.New("unknown pattern", name)
	} else if res, e := pattern.NewResults(run, pat.Name, pat.Return, aff, pat.Labels, args); e != nil {
		err = e
	} else if oldScope, e := run.ReplaceScope(res, true); e != nil {
		err = e
	} else {
		var allFlags rt.Flags
		if rules, e := run.GetRules(pat.Name, "", &allFlags); e != nil {
			err = e
		} else if e := res.ApplyRules(run, rules, allFlags); e != nil {
			err = e
		} else {
			ret, err = res.GetResult()
		}
		// only init can return an error
		run.ReplaceScope(oldScope, false)
	}
	return
}
