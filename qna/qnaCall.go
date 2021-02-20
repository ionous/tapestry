package qna

import (
	"strings"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

// normalize optional arguments
func (run *Runner) Call(pat string, aff affine.Affinity, args []rt.Arg) (ret g.Value, err error) {
	name := lang.Breakcase(pat)
	var labels, result string // fix? consider a cache for this info?
	if e := run.fields.patternOf.QueryRow(name).Scan(&name, &labels, &result); e != nil {
		err = errutil.New("error while querying", pat, e)
	} else {
		if k, e := run.GetKindByName(name); e != nil {
			err = e
		} else {
			rec := k.NewRecord()
			// args run in the scope of their parent context
			// they write to the record that will become the new context
			labels := strings.Split(labels, ",") //
			if e := pattern.DetermineArgs(run, rec, labels, args); e != nil {
				err = e
			} else {
				// locals can ( and often do ) read arguments.
				watcher := pattern.NewResults(rec, result)
				oldScope := run.ReplaceScope(watcher)
				if e := run.initializeLocals(rec); e != nil {
					err = e
				} else {
					var allFlags rt.Flags
					if rules, e := run.GetRules(k.Name(), &allFlags); e != nil {
						err = e
					} else if e := watcher.ApplyRules(run, rules, allFlags); e != nil {
						err = e
					} else {
						ret, err = safe.UnpackResult(rec, result, aff)
					}
				}
				run.ReplaceScope(oldScope)
			}
		}
	}
	if err != nil {
		err = errutil.New("error calling", pat, err)
	}
	return
}

// by now the initializers for the kind will have been cached....
func (run *Runner) initializeLocals(rec *g.Record) (err error) {
	k := rec.Kind()
	if qk, ok := run.qnaKinds.kinds[k.Name()]; !ok {
		err = errutil.New("unknown kind", k.Name())
	} else {
		// run all the initializers
		for i, init := range qk.init {
			if init != nil { // not every field necessarily has an initializer
				if v, e := init.GetAssignedValue(run); e != nil {
					err = errutil.New("error determining local", k.Name(), k.Field(i).Name, e)
					break
				} else if e := rec.SetIndexedField(i, v); e != nil {
					err = errutil.New("error setting local", k.Name(), k.Field(i).Name, e)
					break
				}
			}
		}
	}
	return
}
