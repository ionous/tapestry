package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// note: this is mirrored/mimicked in package testpat
// note: in order to generate appropriate defaults ( ex. a record of the right type )
// can return both a both meaningful value *and* an error
func (run *Runner) Call(name string, aff affine.Affinity, keys []string, vals []g.Value) (ret g.Value, err error) {
	// create a container to hold results of args, locals, and the pending return value
	if kind, e := run.getKind(name); e != nil {
		err = e
	} else {
		if !kind.Implements(kindsOf.Record.String()) {
			// note: this doesnt positively affirm kindsOf.Pattern:
			// some tests use golang structs as faux patterns.
			if v, e := run.call(kind, aff, keys, vals); e != nil {
				err = errutil.Fmt("%w calling %q", e, name)
			} else {
				ret = v
			}
		} else {
			// Call is being used as a side effect to initialize records
			if aff != affine.Record {
				err = errutil.Fmt("attempting to call a record %q with affine %q", name, aff)
			} else if rec, e := safe.FillRecord(run, kind.NewRecord(), keys, vals); e != nil {
				err = e
			} else {
				ret = g.RecordOf(rec)
			}
		}
	}
	return
}

func (run *Runner) call(kind cachedKind, aff affine.Affinity, keys []string, vals []g.Value) (ret g.Value, err error) {
	name := kind.Name()
	if rec, e := safe.FillRecord(run, kind.NewRecord(), keys, vals); e != nil {
		err = e
	} else if res, e := pattern.NewResults(run, rec, aff); e != nil {
		err = e
	} else {
		var flags rt.Flags
		oldScope := run.Stack.ReplaceScope(res)
		run.currentPatterns.startedPattern(name)
		if e := kind.initializeRecord(run, rec); e != nil {
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
				err = errutil.Fmt("%w computing %s", rt.NoResult, aff)
			}
		}
		run.currentPatterns.stoppedPattern(name)
		run.Stack.ReplaceScope(oldScope)
	}
	return
}
