package scope

import (
	"git.sr.ht/~ionous/iffy/object"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

type TargetValue struct {
	Target string
	Value  g.Value
}

func (k *TargetValue) GetField(target, field string) (ret g.Value, err error) {
	if target != object.Variables || field != k.Target {
		err = g.UnknownField{target, field}
	} else {
		ret = k.Value
	}
	return
}

// SetField writes the value of 'v' into the value at 'name'.
func (k *TargetValue) SetField(target, field string, v g.Value) (err error) {
	if target != object.Variables || field != k.Target {
		err = g.UnknownField{target, field}
	} else if g.AreEqualTypes(k.Value, v) {
		k.Value = v
	} else {
		err = errutil.New("mismatched types")
	}
	return
}
