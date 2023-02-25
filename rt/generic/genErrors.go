package generic

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/rt/meta"
	"github.com/ionous/errutil"
)

type Overflow struct {
	Index, Bounds int
}

func (e Overflow) Error() string {
	return errutil.Sprint(e.Index, "above range", e.Bounds)
}

type Underflow struct {
	Index, Bounds int
}

func (e Underflow) Error() string {
	return errutil.Sprint(e.Index, "below range", e.Bounds)
}

// error for GetField, SetField
type Unknown struct {
	Target, Field string
}

func (e Unknown) IsUnknownField() bool {
	tgt := e.Target
	return len(tgt) > 0 && tgt != meta.Variables && tgt != meta.ObjectId
}

func (e Unknown) Error() (ret string) {
	switch e.Target {
	case "":
		ret = errutil.Sprintf("unknown name %q", e.Field)
	case meta.Variables:
		ret = errutil.Sprintf("unknown variable %q", e.Field)
	case meta.ObjectId:
		ret = errutil.Sprintf("unknown object %q", e.Target)
	default:
		ret = errutil.Sprintf(`unknown field "%s.%s"`, e.Target, e.Field)
	}
	return
}

func UnknownName(name string) error {
	return Unknown{Field: name}
}

func UnknownVariable(v string) error {
	return Unknown{Target: meta.Variables, Field: v}
}

func UnknownObject(o string) error {
	return Unknown{Target: meta.ObjectId, Field: o}
}

func UnknownField(target, field string) error {
	return Unknown{Target: target, Field: field}
}

func IsUnknown(e error) bool {
	var u Unknown
	return errors.As(e, &u)
}
