package generic

import (
	"github.com/ionous/errutil"
)

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
	} else if len(e.Field) == 0 {
		ret = errutil.Sprintf("unknown object %q", e.Target)
	} else {
		ret = errutil.Sprintf(`unknown field "%s.%s"`, e.Target, e.Field)
	}
	return
}

func UnknownVariable(v string) error {
	return Unknown{Field: v}
}

func UnknownObject(o string) error {
	return Unknown{Target: o}
}

func UnknownField(target, field string) error {
	return Unknown{Target: target, Field: field}
}
