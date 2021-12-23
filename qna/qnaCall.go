package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/pattern"
	"github.com/ionous/errutil"
)

// note: this is mirrored/mimicked in package testpat
func (run *Runner) Call(pat string, aff affine.Affinity, args []rt.Arg) (ret g.Value, err error) {
	name := lang.Underscore(pat) // FIX: why are people calling this with untransformed names
	if pl, e := run.qdb.PatternLabels(name); e != nil {
		err = e
	} else if rec, e := pattern.NewRecord(run, name, pl.Labels, args); e != nil {
		err = e
	} else {
		// locals can ( and often do ) read arguments ( which can invoke sub-patterns )
		run.currentPatterns.startedPattern(name)
		results := pattern.NewResults(rec, pl.Result, aff)
		if oldScope, e := run.ReplaceScope(results, true); e != nil {
			err = e
		} else {
			var flags rt.Flags
			if rules, e := run.GetRules(name, "", &flags); e != nil {
				err = e
			} else if e := results.ApplyRules(run, rules, flags); e != nil {
				err = e
			} else if v, e := results.GetResult(); e != nil {
				err = e
			} else {
				// breaks precedence to return a value and an error
				// in order to generate appropriate default returns ( ex. a record of the right type )
				// while still informing the caller of lack of pattern decision in a concise manner.
				ret = v
				if !results.ComputedResult() {
					err = rt.NoResult{}
				}
			}
			// only init can return an error
			run.ReplaceScope(oldScope, false)
			run.currentPatterns.stoppedPattern(name)
		}
	}
	if err != nil {
		// err = errutil.Fmt("%w calling %s with %v", err, pat, g.RecordToValue(rec))
		err = errutil.Fmt("%w calling %s", err, pat)
	}
	return
}
