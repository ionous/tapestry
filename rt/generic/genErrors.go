package generic

import (
	"git.sr.ht/~ionous/iffy/object"
	"github.com/ionous/errutil"
)

// error constant for iterators
const NothingObject errutil.Error = "object is nothing"

type Overflow struct {
	Index, Bounds int
}

type Underflow struct {
	Index, Bounds int
}

// error for GetField, SetField
type Unknown struct {
	Target, Field string
}

func (e Underflow) Error() string {
	return errutil.Sprint(e.Index, "below range", e.Bounds)
}

func (e Overflow) Error() string {
	return errutil.Sprint(e.Index, "above range", e.Bounds)
}

func (e Unknown) Error() (ret string) {
	if len(e.Target) == 0 {
		ret = errutil.Sprintf("unknown variable %q", e.Field)
	} else if e.Target == object.Value {
		ret = errutil.Sprintf("unknown object %q", e.Field)
	} else {
		ret = errutil.Sprintf(`unknown field "%s.%s"`, e.Target, e.Field)
	}
	return
}

func UnknownVariable(v string) error {
	return Unknown{Field: v}
}

func UnknownObject(o string) error {
	return Unknown{Target: object.Value, Field: o}
}

func UnknownField(target, field string) error {
	return Unknown{Target: target, Field: field}
}
