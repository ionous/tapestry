package pattern

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func GetResultField(run rt.Runtime, k *g.Kind) (ret int, err error) {
	patternName := k.Name()
	if labels, e := run.GetField(meta.PatternLabels, patternName); e != nil {
		err = e
	} else {
		labels := labels.Strings()
		if cnt := len(labels); cnt == 0 || labels[cnt-1] == "" {
			ret = -1
		} else {
			fieldName := labels[cnt-1]
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
	rec       *g.Record
	field     int
	hasResult bool
}

func (res *Result) GetResult(aff affine.Affinity) (ret g.Value, err error) {
	rec, field, okay := res.rec, res.field, res.hasResult
	if field < 0 {
		// no result field, but we still might be checking for whether it had any matching rules.
		if aff == affine.Bool {
			ret = g.BoolOf(okay)
		} else if len(aff) != 0 {
			err = errutil.Fmt("%w; caller expected %s", rt.NoResult, aff)
		}
	} else {
		// get the value *or* a default.
		if v, e := rec.GetIndexedField(field); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else if e := safe.Check(v, aff); e != nil {
			err = errutil.New("error trying to get return value", e)
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
			// we can return both an error *and* a (default) value.
			if !okay {
				err = rt.NoResult
			}
			ret = v
		}
	}
	return
}
