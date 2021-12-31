package safe

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
)

// used when converting values to fields that might require objects.
// if the target field (ex. a pattern local) requires text of a certain type
// and the incoming value is untyped, try to convert it.
func AutoConvert(run rt.Runtime, ft g.Field, val g.Value) (ret g.Value, err error) {
	ret = val // provisionally.
	// assigning to typed text?
	if ft.Affinity == affine.Text && len(ft.Type) > 0 {
		// and we have some valid untyped text.
		if val.Affinity() == affine.Text && val.Len() > 0 && len(val.Type()) == 0 {
			ret, err = run.GetField(meta.ObjectId, val.String())
		}
	}
	return
}
