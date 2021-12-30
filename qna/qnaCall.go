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
	name := lang.Underscore(pat) // fix: why are people calling this with untransformed names
	// note: locals can ( and often do ) read arguments ( which can invoke sub-patterns )
	if pl, e := run.qdb.PatternLabels(name); e != nil {
		err = e
	} else if res, e := pattern.NewResults(run, name, pl.Result, aff, pl.Labels, args); e != nil {
		err = e
	} else {
		run.currentPatterns.startedPattern(name)
		// note: local Init happens inside of ReplaceScope :/
		if oldScope, e := run.ReplaceScope(res, true); e != nil {
			err = e
		} else {
			var flags rt.Flags
			if rules, e := run.GetRules(name, "", &flags); e != nil {
				err = e
			} else if e := res.ApplyRules(run, rules, flags); e != nil {
				err = e
			} else if v, e := res.GetResult(); e != nil {
				err = e
			} else {
				// breaks precedence to return a value and an error
				// in order to generate appropriate default returns ( ex. a record of the right type )
				// while still informing the caller of lack of pattern decision in a concise manner.
				ret = v
				if !res.ComputedResult() {
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
