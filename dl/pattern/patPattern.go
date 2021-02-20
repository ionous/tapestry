package pattern

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

type Pattern struct {
	Name   string
	Return string          // name of return field; empty if none ( could be an index but slightly safer this way )
	Labels []string        // one label for every parameter
	Locals []rt.Assignment // usually equal to the number of locals; or nil for testing.
	Fields []g.Field       // flat list of params and locals and an optional return
	Rules  []rt.Rule
}

func (pat *Pattern) Run(run rt.Runtime, aff affine.Affinity, args []rt.Arg) (ret g.Value, err error) {
	// create a container to hold results of args, locals, and the pending return value
	if k, e := run.GetKindByName(pat.Name); e != nil {
		err = e
	} else {
		rec := k.NewRecord()
		// args run in the scope of their parent context
		// they write to the record that will become the new context
		if e := DetermineArgs(run, rec, pat.Labels, args); e != nil {
			err = e
		} else {
			// initializers ( and the pattern itself ) run in the scope of the pattern
			// ( with access to all locals and args)
			watcher := NewResults(rec, pat.Return)
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
	return
}
