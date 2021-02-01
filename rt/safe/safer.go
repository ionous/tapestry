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
