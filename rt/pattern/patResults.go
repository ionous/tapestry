package pattern

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

// track the number of times a particular field gets successfully written to.
// implements rt.Scope
type Results struct {
	rec                  *g.Record
	resultField          string
	expectedAff          affine.Affinity
	resultSets, ranCount int
	ranFlags             rt.Flags
}

// fix; itd be nice at least to cleanout expectedAff if possible.
func NewResults(run rt.Runtime, name string, res string, aff affine.Affinity, parts []string, args []rt.Arg) (ret *Results, err error) {
	if rec, e := newPattern(run, name, parts, args); e != nil {
		err = e
	} else {
		ret = &Results{
			rec:         rec,
			resultField: res,
			expectedAff: aff,
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
		rw.resultSets++
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
	return rw.resultSets > 0 || (len(rw.expectedAff) == 0 && rw.ranCount > 0)
}

func (rw *Results) reset() {
	rw.resultSets = 0
	rw.ranCount = 0
	rw.ranFlags = 0
}

// GetResult returns a default value if none was computed.
func (rw *Results) GetResult() (ret g.Value, err error) {
	field, aff := rw.resultField, rw.expectedAff
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
	if flags := rule.Flags(); (allow&flags != 0) || (flags&rt.Filter != 0) { // Filter means run always
		if ok, e := safe.GetOptionalBool(run, rule.Filter, true); e != nil {
			err = e
		} else if ok.Bool() && allow&flags != 0 {
			if e := safe.Run(run, rule.Execute); e != nil {
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
