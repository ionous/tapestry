package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"github.com/ionous/errutil"
)

// note: this is mirrored/mimicked in package testpat
// note: in order to generate appropriate defaults ( ex. a record of the right type )
// can return both a both meaningful value *and* an error
func (run *Runner) Call(rec *g.Record, aff affine.Affinity) (ret g.Value, err error) {
	name := rec.Kind().Name()
	if res, e := pattern.NewResults(run, rec, aff); e != nil {
		err = e
	} else {
		oldScope := run.Stack.ReplaceScope(res)
		if cached, e := run.getKindOf(name, kindsOf.Pattern.String()); e != nil {
			err = e
		} else {
			var flags rt.Flags
			run.currentPatterns.startedPattern(name)
			if e := cached.initializeRecord(run, rec); e != nil {
				err = e
			} else if rules, e := run.GetRules(name, "", &flags); e != nil {
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
					err = errutil.Fmt("%w calling %s pattern %q", rt.NoResult, aff, name)
				}
			}
			run.currentPatterns.stoppedPattern(name)
		}
		run.Stack.ReplaceScope(oldScope)
	}
	if err != nil {
		// err = errutil.Fmt("%w calling %s with %v", err, pat, g.RecordToValue(rec))
		err = errutil.Fmt("%w calling %s", err, name)
	}
	return
}
