package testpat

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
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
		// create a container to hold results of args, locals, and the pending return value
		if k, e := run.GetKindByName(pat.Name); e != nil {
			err = e
		} else {
			rec := k.NewRecord()
			// args run in the scope of their parent context
			// they write to the record that will become the new context
			if e := pattern.DetermineArgs(run, rec, pat.Labels, args); e != nil {
				err = e
			} else {
				// initializers ( and the pattern itself ) run in the scope of the pattern
				// ( with access to all locals and args)
				watcher := pattern.NewResults(rec, pat.Return)
				oldScope := run.ReplaceScope(watcher)
				var allFlags rt.Flags
				if rules, e := run.GetRules(pat.Name, &allFlags); e != nil {
					err = e
				} else if e := watcher.ApplyRules(run, rules, allFlags); e != nil {
					err = e
				} else {
					ret, err = safe.UnpackResult(rec, pat.Return, aff)
				}
				run.ReplaceScope(oldScope)
			}
		}
	}
	return
}
