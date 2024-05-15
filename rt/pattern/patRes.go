package pattern

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func GetResultField(run rt.Runtime, k *rt.Kind) (ret int, err error) {
	ret = -1 // provisionally
	patternName := k.Name()
	if labels, e := run.GetField(meta.PatternLabels, patternName); e != nil {
		err = e
	} else {
		labels := labels.Strings()
		if end := len(labels) - 1; end >= 0 && len(labels[end]) > 0 {
			fieldName := labels[end]
			if i := k.FieldIndex(fieldName); i < 0 {
				err = errutil.New("couldn't find return field %q in kind %q", fieldName, patternName)
			} else {
				ret = i
			}
		}
	}
	return
}

type Result struct {
	rec         *rt.Record
	resultField int
	hasResult   bool
}

func (res *Result) GetResult(run rt.Runtime, aff affine.Affinity) (ret rt.Value, err error) {
	rec, resultField, hasResult := res.rec, res.resultField, res.hasResult
	if resultField < 0 {
		// no result field, but we still might be checking for whether it had any matching rules.
		if aff == affine.Bool {
			ret = rt.BoolOf(hasResult)
		} else if len(aff) != 0 { // fix: why is this even being run in the first place?
			err = errutil.Fmt("no result when caller expected %s", aff)
		}
	} else {
		// get the value *or* a default.
		if v, e := rec.GetIndexedField(resultField); e != nil {
			err = errutil.New("error getting result", e)
		} else if v, e := safe.ConvertValue(run, v, aff); e != nil {
			err = errutil.New("error checking result", e)
		} else {
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
			ret = v
		}
	}
	return
}
