package rt

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt/meta"
)

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
		ret = fmt.Sprintf("unknown name %q", e.Field)
	case meta.ObjectId:
		ret = fmt.Sprintf("unknown object %q", e.Field)
	case meta.Response:
		ret = fmt.Sprintf("unknown response %q", e.Field)
	case meta.Variables:
		ret = fmt.Sprintf("unknown variable %q", e.Field)
	default:
		ret = fmt.Sprintf(`unknown field "%s.%s"`, e.Target, e.Field)
	}
	return
}

func UnknownName(name string) error {
	return Unknown{Field: name}
}

func UnknownResponse(v string) error {
	return Unknown{Target: meta.Response, Field: v}
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
