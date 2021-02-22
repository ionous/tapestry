package safe

import (
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

func GetAssignedValue(run rt.Runtime, a rt.Assignment) (ret g.Value, err error) {
	if a == nil {
		err = MissingEval("assignment")
	} else {
		ret, err = a.GetAssignedValue(run)
	}
	return
}
