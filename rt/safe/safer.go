package safe

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

func Check(v g.Value, want affine.Affinity) (err error) {
	if va := v.Affinity(); len(want) > 0 && want != va {
		err = errutil.Fmt("have %q, wanted %q", va, want)
	}
	return
}

// resolve a requested variable name into a value of the desired affinity.
func CheckVariable(run rt.Runtime, n string, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := run.GetField(object.Variables, n); e != nil {
		err = e
	} else if e := Check(v, aff); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}

func Unpack(src g.Value, field string, aff affine.Affinity) (ret g.Value, err error) {
	if !affine.HasFields(src.Affinity()) {
		err = errutil.New("Value", src, "doesn't have fields")
	} else if v, e := src.FieldByName(field); e != nil {
		err = e
	} else if e := Check(v, aff); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}

// fix: still trying to figure out where this should live
// maybe just merge with regular unpack? ( though FieldByName copies and GetNamedField does not )
func UnpackResult(src *g.Record, field string, aff affine.Affinity) (ret g.Value, err error) {
	if len(field) > 0 {
		// get the value and check its result
		if v, e := src.GetNamedField(field); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else if e := Check(v, aff); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else if len(aff) == 0 {
			// the caller expects nothing but we have a return value.
			if v.Affinity() == affine.Text {
				HackTillTemplatesCanEvaluatePatternTypes = v.String()
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

var HackTillTemplatesCanEvaluatePatternTypes string
