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
	resultField string
	resultAff   affine.Affinity
	sets        int
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
		rw.sets++
		// we could also store the last value set,
		// and use that for our result --
		// but if nothing is set, it feels better to use the default record value production
		// generated by record.GetNamedField
	}
	return
}

func (rw *Results) Compute(run rt.Runtime, rules []rt.Rule, allFlags rt.Flags) (ret g.Value, err error) {
	if e := rw.ApplyRules(run, rules, allFlags); e != nil {
		err = e
	} else {
		ret, err = rw.GetResult()
	}
	return
}

func (rw *Results) HasResults() bool {
	return rw.sets > 0
}

func (rw *Results) GetResult() (ret g.Value, err error) {
	rec := rw.Scope.(g.Value).Record()
	field, aff := rw.resultField, rw.resultAff
	if len(field) > 0 {
		// get the value and check its result
		if v, e := rec.GetNamedField(field); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else if e := safe.Check(v, aff); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else if len(aff) == 0 {
			// the caller expects nothing but we have a return value.
			if v.Affinity() == affine.Text {
				safe.HackTillTemplatesCanEvaluatePatternTypes = v.String()
			}
			// other than passing data back to templates in a hack...
			// we dont treat this as an error -- we allow patterns to be run for side effects.
		} else {
			ret = v
		}
	} else if len(aff) != 0 {
		err = errutil.New("caller expected", aff, "returned nothing")
	}
	return
}

// RunWithScope - note: assumes whatever scope is needed to run the pattern has already been setup.
func (rw *Results) ApplyRules(run rt.Runtime, rules []rt.Rule, allFlags rt.Flags) (err error) {
	sets := rw.sets
	for i, cnt := 0, len(rules); i < cnt && allFlags != 0; i++ {
		if ranFlag, e := ApplyRule(run, rules[i], allFlags); e != nil {
			err = e
		} else if ranFlag != 0 {
			didSomething := (rw.sets > sets)
			sets = rw.sets
			// if we ran a prefix or a post fix rule and it did something, we are done.
			if didSomething && ranFlag != rt.Infix {
				break
			}
			// otherwise, if an infix rule did something
			// check the other kinds of rules
			// ditto if we dont expect the pattern to return anything:
			// in that case we just want to do the first of each rule type.
			if didSomething || len(rw.resultField) == 0 {
				allFlags = allFlags &^ ranFlag
			}
		}
	}
	return
}

func ApplyRule(run rt.Runtime, rule rt.Rule, allow rt.Flags) (ret rt.Flags, err error) {
	if flags := rule.GetFlags(); allow&flags != 0 {
		if ok, e := safe.GetOptionalBool(run, rule.Filter, true); e != nil {
			err = e
		} else if ok.Bool() { // the rule returns false if it didnt apply
			if e := safe.Run(run, rule.Execute); e != nil {
				err = e
			} else {
				ret = flags
			}
		}
	}
	return
}
