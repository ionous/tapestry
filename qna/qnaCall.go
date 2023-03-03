package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"github.com/ionous/errutil"
)

// note: this is mirrored/mimicked in package testpat
// note: in order to generate appropriate defaults ( ex. a record of the right type )
// can return both a both meaningful value *and* an error
func (run *Runner) Call(pat *g.Record, aff affine.Affinity) (ret g.Value, err error) {
	name := pat.Kind().Name()
	if res, e := pattern.NewResults(run, pat, aff); e != nil {
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
				// warning: in order to generate appropriate defaults ( ex. a record of the right type )
				// while still informing the caller of lack of pattern decision in a concise manner
				// can return both a valid value and an error
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
		err = errutil.Fmt("%w calling %s", err, name)
	}
	return
}
