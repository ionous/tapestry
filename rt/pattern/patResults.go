package pattern

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

// track the number of times a particular field gets successfully written to.
type Results struct {
	rt.Scope
	resultField          string
	resultAff            affine.Affinity
	resultSets, ranCount int
	ranFlags             rt.Flags
}

func NewResults(rec *g.Record, res string, aff affine.Affinity) *Results {
	return &Results{
		Scope:       g.RecordOf(rec),
		resultField: res,
		resultAff:   aff,
	}
}

func (rw *Results) SetFieldByName(field string, val g.Value) (err error) {
	if e := rw.Scope.SetFieldByName(field, val); e != nil {
		err = e
	} else if field == rw.resultField {
		rw.resultSets++
		// we could also store the last value set,
		// and use that for our result --
		// but if nothing is set, it feels better to use the default record value production
		// generated by record.GetNamedField
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
	return rw.resultSets > 0 || (len(rw.resultAff) == 0 && rw.ranCount > 0)
}

func (rw *Results) reset() {
	rw.resultSets = 0
	rw.ranCount = 0
	rw.ranFlags = 0
}

// GetResult returns a default value if none was computed.
func (rw *Results) GetResult() (ret g.Value, err error) {
	field, aff := rw.resultField, rw.resultAff
	if len(field) == 0 {
		// no result field, but we still might be checking for whether it had any matching rules.
		if aff == affine.Bool {
			// should it be any rule? just the infix rule?
			// probably always using a return flag would be best.
			ret = g.BoolOf(rw.ranFlags&rt.Infix != 0)
		} else if len(aff) != 0 {
			err = errutil.New("caller expected", aff, "returned nothing")
		}
	} else {
		// get the value and check its result
		rec := rw.Scope.(g.Value).Record()
		if v, e := rec.GetNamedField(field); e != nil {
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
				safe.HackTillTemplatesCanEvaluatePatternTypes = v.String()
			}
		}
	}
	return
}

// add a member apply rule so its easier to call singly
// -- maybe allFlags is a pointer so you can manage it better?
// outside youll have to set value scope OR AddScope value target.
// ( the latter might be easier for now )
// call it "current noun" and allow "current kind" as well
// --

// ApplyRules - note: assumes whatever scope is needed to run the pattern has already been setup.
func (rw *Results) ApplyRules(run rt.Runtime, rules []rt.Rule, flags rt.Flags) (err error) {
	for _, rule := range rules {
		if next, e := rw.ApplyRule(run, rule, flags); e != nil || next == 0 {
			err = e
			break
		} else {
			flags = next
		}
	}
	return
}

// ApplyRule - note: assumes whatever scope is needed to run the pattern has already been setup.
// returns remaining flags.
func (rw *Results) ApplyRule(run rt.Runtime, rule rt.Rule, flags rt.Flags) (ret rt.Flags, err error) {
	resultSets := rw.resultSets // check if rule changes this.
	if ranFlag, e := ApplyRule(run, rule, flags); e != nil {
		err = errutil.New(e, "while applying", rule.Name)
	} else {
		var didSomething bool
		if ranFlag != 0 {
			rw.ranCount++
			rw.ranFlags |= ranFlag
			// did the rule return a value ( or did it run and doesn't expect an explicit return )
			didSomething = len(rw.resultField) == 0 || rw.resultSets > resultSets
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
func ApplyRule(run rt.Runtime, rule rt.Rule, allow rt.Flags) (ret rt.Flags, err error) {
	if flags := rule.GetFlags(); allow&flags != 0 {
		if ok, e := safe.GetOptionalBool(run, rule.Filter, true); e != nil {
			err = e
		} else if ok.Bool() {
			if e := safe.Run(run, rule.Execute); e != nil {
				err = e
			} else {
				ret = flags
			}
		}
	}
	return
}
