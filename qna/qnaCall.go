package qna

import (
	"database/sql"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/pattern"
	"github.com/ionous/errutil"
)

// note: this is mirrored/mimicked in package testpat
func (run *Runner) Call(pat string, aff affine.Affinity, args []rt.Arg) (ret g.Value, err error) {
	name := lang.Breakcase(pat) // gets replaced with the actual name by query
	var labels, result string   // fix? consider a cache for this info?
	var rec *g.Record
	if e := run.fields.patternOf.QueryRow(name).Scan(&name, &labels, &result); e != nil {
		if e == sql.ErrNoRows {
			err = errutil.Fmt("couldn't find the pattern named %q", name)
		} else {
			err = e
		}
	} else if rec, e = pattern.NewRecord(run, name, labels, args); e != nil {
		err = e
	} else {
		// locals can ( and often do ) read arguments ( which can invoke sub-patterns )
		run.currentPatterns.startedPattern(name)
		results := pattern.NewResults(rec, result, aff)
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
