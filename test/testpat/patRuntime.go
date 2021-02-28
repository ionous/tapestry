package testpat

import (
	"strings"

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
	} else {
		ls := strings.Join(pat.Labels, ",")
		if rec, e := pattern.NewRecord(run, pat.Name, ls, args); e != nil {
			err = e
		} else {
			results := pattern.NewResults(rec, pat.Return, aff)
			if oldScope, e := run.ReplaceScope(results, true); e != nil {
				err = e
			} else {
				var allFlags rt.Flags
				if rules, e := run.GetRules(pat.Name, "", &allFlags); e != nil {
					err = e
				} else {
					ret, err = results.Compute(run, rules, allFlags)
				}
				// only init can return an error
				run.ReplaceScope(oldScope, false)
			}
		}
	}
	return
}
