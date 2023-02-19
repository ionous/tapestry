package pattern

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// track the number of times a particular field gets successfully written to.
// implements rt.Scope
type Results struct {
	rec         *g.Record
	resultField string          // variable name for the result of the pattern (if any)
	expectedAff affine.Affinity // the type of the result
	resultCount int             // number of times the result field was written to
	currRules   int             // number of rules that were run
	currFlags   rt.Flags        // which phases have run, still have to run?
}

// fix. itd be nice to remove expectedAff if possible. ( need to fx HackTillTemplatesCanEvaluatePatternTypes i think )
func NewResults(run rt.Runtime, rec *g.Record, aff affine.Affinity) (ret *Results, err error) {
	// tbd: is the result field really last?
	// fix: if every pattern returned a result (ex. if execute returned a hidden bool or c-like int )
	// then we could avoid the label look up, and always use the last indexed field
	// currently, there might be a blank label "" for the return, which means no indexed field
	// we also cant use "aff" because the calling context might not be right ( re Hack, and hints )
	if rec == nil {
		err = errutil.New("internal error: nil pattern passed to pattern results")
	} else {
		name := rec.Kind().Name()
		if labels, e := run.GetField(meta.PatternLabels, name); e != nil {
			err = e
		} else {
			labels := labels.Strings()
			if last := len(labels) - 1; last < 0 {
				err = errutil.New("pattern has unexpectedly few labels", name)
			} else {
				ret = &Results{
					rec:         rec,
					resultField: labels[last],
					expectedAff: aff,
				}
			}
		}
	}
	return
}

func (rw *Results) Record() *g.Record {
	return rw.rec
}

// implements rt.Scope
func (rw *Results) FieldByName(field string) (g.Value, error) {
	return rw.rec.GetNamedField(field)
}

// implements rt.Scope
func (rw *Results) SetFieldByName(field string, val g.Value) (err error) {
	// fix: what happens if you set a trait and the aspect should have been the return value
	// it would be hidden inside SetFieldByName -- and we'd never see it.
	if e := rw.rec.SetNamedField(field, val); e != nil {
		err = e
	} else if field == rw.resultField {
		rw.resultCount++
	}
	return
}

// returns false if the computed result was false ( and on error )
func (rw *Results) GetContinuation() (okay bool, err error) {
	if !rw.ComputedResult() {
		okay = true
	} else if res, e := rw.GetResult(); e != nil {
		err = e
	} else if !res.Bool() {
		okay = false
	} else {
		rw.reset() // reset the internals so we can keep using the results object
		okay = true
	}
	return
}

// ComputedResult returns whether an explicit result was set.
func (rw *Results) ComputedResult() bool {
	// did it compute a result; or -- if it wasnt expecting a result -- did at least something happen?
	// fix: right now aff is coming from the caller. that seems wrong.
	// why cant the caller just safe.Check the result is what they want.
	return rw.resultCount > 0 || (len(rw.expectedAff) == 0 && rw.currRules > 0)
}

func (rw *Results) reset() {
	rw.resultCount = 0
	rw.currRules = 0
	rw.currFlags = 0
}

// GetResult returns a default value if none was computed.
func (rw *Results) GetResult() (ret g.Value, err error) {
	field, aff := rw.resultField, rw.expectedAff
	if len(field) == 0 {
		// no result field, but we still might be checking for whether it had any matching rules.
		if aff == affine.Bool {
			// should it be any rule? just the infix rule?
			// probably always using a return flag would be best.
			ret = g.BoolOf(rw.currFlags&rt.Infix != 0)
		} else if len(aff) != 0 {
			err = errutil.New("caller expected", aff, "returned nothing")
		}
	} else {
		// get the value and check its result
		if v, e := rw.rec.GetNamedField(field); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else if e := safe.Check(v, aff); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else {
			ret = v
			// the caller expects nothing but we have a return value.
			// other than passing data back to templates in a hack...
			// we dont treat this as an error --
			// a) we allow patterns to be run for side effects, and:
			// b) "trying" doesnt know the affinity of the value at the time of the call.
			if len(aff) == 0 && v.Affinity() == affine.Text {
				safe.HackTillTemplatesCanEvaluatePatternTypes = v
			} else {
				safe.HackTillTemplatesCanEvaluatePatternTypes = nil
			}
		}
	}
	return
}

// ApplyRules - note: assumes whatever scope is needed to run the pattern has already been setup.
func (rw *Results) ApplyRules(run rt.Runtime, rules []rt.Rule, flags rt.Flags) (err error) {
	for _, rule := range rules {
		// end if there are no flags left, and we didn't want to filter everything.
		if next, e := rw.ApplyRule(run, rule, flags); e != nil || (next == 0 && flags&rt.Filter == 0) {
			err = e
			break
		} else {
			flags = next | (flags & rt.Filter) // if filter was set, keep it set.
		}
	}
	return
}

// ApplyRule - assumes whatever scope is needed to run the pattern has already been setup.
// returns remaining flags.
func (rw *Results) ApplyRule(run rt.Runtime, rule rt.Rule, flags rt.Flags) (ret rt.Flags, err error) {
	resultCount := rw.resultCount // check if rule changes this.
	if ranFlag, e := apply(run, rule, flags); e != nil {
		err = errutil.New(e, "while applying", rule.Name)
	} else {
		var didSomething bool
		if ranFlag != 0 {
			rw.currRules++
			rw.currFlags |= ranFlag
			// did the rule return a value ( or did it run and doesn't expect an explicit return )
			didSomething = len(rw.resultField) == 0 || rw.resultCount > resultCount
		}
		//
		if !didSomething {
			ret = flags // no? keep trying rules of this type.
		} else if ranFlag == rt.Infix {
			// for prefix and postfix rules, we are completely done: return 0.
			// for infix rules, just stop running infix type
			ret = flags &^ ranFlag
		}
	}
	return
}

// return the flags of the rule if it ran; even if it didnt return anything.
func apply(run rt.Runtime, rule rt.Rule, allow rt.Flags) (ret rt.Flags, err error) {
	// get the rule's flags and see whether we should run the rule
	// rt.Filter is a flag for filters that need to always evalute (ex. counters)
	if flags := rule.Flags(); (allow&flags != 0) || (flags&rt.Filter != 0) {
		// actually run the filter:
		if ok, e := safe.GetOptionalBool(run, rule.Filter, true); e != nil {
			err = e
		} else if ok.Bool() && allow&flags != 0 { // check whether the filter succeeded and whether the rule is actually supposed to run.
			if e := safe.RunAll(run, rule.Execute); e != nil {
				err = e
			} else {
				// don't let the "run always" filter flag show through
				// only let through if we actually executed something.
				ret = flags & ^rt.Filter
			}
		}
	}
	return
}
