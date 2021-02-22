package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

// note: this is mirrored/mimicked in package testpat
func (run *Runner) Call(pat string, aff affine.Affinity, args []rt.Arg) (ret g.Value, err error) {
	name := lang.Breakcase(pat)
	var labels, result string // fix? consider a cache for this info?
	if e := run.fields.patternOf.QueryRow(name).Scan(&name, &labels, &result); e != nil {
		err = errutil.New("error while querying", pat, e)
	} else if rec, e := pattern.NewRecord(run, name, labels, args); e != nil {
		err = e
	} else {
		// locals can ( and often do ) read arguments.
		results := pattern.NewResults(rec, result, aff)
		if oldScope, e := run.ReplaceScope(results, true); e != nil {
			err = e
		} else {
			var allFlags rt.Flags
			if rules, e := run.GetRules(name, "", &allFlags); e != nil {
				err = e
			} else {
				ret, err = results.Compute(run, rules, allFlags)
			}
			// only init can return an error
			run.ReplaceScope(oldScope, false)
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
