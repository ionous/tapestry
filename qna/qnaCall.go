package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

// normalize optional arguments
func (run *Runner) Call(name string, aff affine.Affinity, args []rt.Arg) (ret g.Value, err error) {
	pat := new(pattern.Pattern) // FIX
	if e := run.GetEvalByName(name, pat); e != nil {
		err = e
	} else {
		if k, e := run.GetKindByName(name); e != nil {
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
				// locals ( by definition ) write to the record context
				if e := run.initializeLocals(pat, rec); e != nil {
					err = e
				} else {
					var allFlags rt.Flags
					if rules, e := run.GetRules(k.Name(), &allFlags); e != nil {
						err = e
					} else if e := watcher.ApplyRules(run, rules, allFlags); e != nil {
						err = e
					} else {
						ret, err = safe.UnpackResult(rec, pat.Return, aff)
					}
				}
				run.ReplaceScope(oldScope)
			}
		}
	}
	if err != nil {
		err = errutil.New("error calling", name, err)
	}
	return
}

func (run *Runner) initializeLocals(pat *pattern.Pattern, rec *g.Record) (err error) {
	lin, fin, lcnt := 0, len(pat.Labels), len(pat.Locals) // locals start after labels
	k := rec.Kind()
	for lin < lcnt {
		if field, init := k.Field(fin), pat.Locals[lin]; init != nil {
			if v, e := init.GetAssignedValue(run); e != nil {
				err = errutil.New(pat.Name, "error determining local", lin, field.Name, e)
				break
			} else if e := rec.SetIndexedField(fin, v); e != nil {
				err = errutil.New(pat.Name, "error setting local", lin, field.Name, e)
				break
			}
		}
		lin++
		fin++
	}
	return
}
