package render

import (
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// RunTest returns an error on failure.
func (op *RenderTemplate) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.GetOptionalText(run, op.Expression, ""); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
