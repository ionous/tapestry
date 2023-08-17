package pattern

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
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
}

// fix. itd be nice to remove expectedAff if possible. ( need to fx HackTillTemplatesCanEvaluatePatternTypes i think )
func NewResults(rec *g.Record, labels []string, aff affine.Affinity) (ret *Results) {
	// we allow nil labels as a shortcut for events;
	var resultField string
	if cnt := len(labels); cnt > 0 {
		resultField = labels[cnt-1]
	}
	return &Results{
		rec:         rec,
		resultField: resultField,
		expectedAff: aff,
	}
}

func (rw *Results) Record() *g.Record {
	return rw.rec
}

// implements rt.Scope
func (rw *Results) FieldByName(field string) (g.Value, error) {
	return rw.rec.GetNamedField(field)
}

// implements rt.Scope
func (rw *Results) SetFieldDirty(field string) (err error) {
	if field == rw.resultField {
		rw.resultCount++
	} else {
		// todo: example, flag object or db for save.
		// for now, simply verify that the field exists.
		_, err = rw.rec.GetNamedField(field)
	}
	return
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
}

// GetResult returns a default value if none was computed.
func (rw *Results) GetResult() (ret g.Value, err error) {
	field, aff := rw.resultField, rw.expectedAff
	if len(field) == 0 {
		// no result field, but we still might be checking for whether it had any matching rules.
		if aff == affine.Bool {
			ret = g.BoolOf(rw.ComputedResult())
		} else if len(aff) != 0 {
			err = errutil.Fmt("%w; caller expected %s", rt.NoResult, aff)
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
func (rw *Results) ApplyRules(run rt.Runtime, rules []rt.Rule) (err error) {
	var skip bool
	for _, rule := range rules {
		// end if there are no flags left, and we didn't want to filter everything.
		if skip, err = rw.ApplyRule(run, rule, skip); err != nil {
			break
		}
	}
	return
}

// ApplyRule - assumes whatever scope is needed to run the pattern has already been setup.
// skip controls whether the rules should be executed or whether the filters need updating.
func (rw *Results) ApplyRule(run rt.Runtime, rule rt.Rule, skip bool) (_ bool, err error) {
	resultCount := rw.resultCount // check if rule changes this.
	if rule.Updates || !skip {
		if ok, e := safe.GetOptionalBool(run, rule.Filter, true); e != nil {
			err = e
		} else if ok.Bool() && !skip {
			if e := safe.RunAll(run, rule.Execute); e != nil {
				err = e
			} else {
				rw.currRules++
				// did the rule return a value ( or did it run and doesn't expect an explicit return )
				skip = len(rw.resultField) == 0 || rw.resultCount > resultCount
			}
		}
	}
	return skip, err
}
