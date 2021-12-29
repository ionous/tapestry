package safe

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/meta"
)

// used when converting values to fields that might require objects.
// if the target field (ex. a pattern local) requires text of a certain type
// and the incoming value is untyped: convert it.
func AutoConvert(run rt.Runtime, ft g.Field, val g.Value) (ret g.Value, err error) {
	ret = val // provisionally.
	if ft.Affinity == affine.Text && len(ft.Type) > 0 {
		if (val.Affinity() == affine.Text && len(val.Type()) == 0) ||
			(val.Affinity() == affine.Object) {
			if str := val.String(); len(str) == 0 {
				ret = g.StringFrom("", ft.Type) // allow setting of blank values
			} else {
				// set indexed field validates the ft.Type and the val.Type match
				// we just have to give it the proper value in the first place.
				ret, err = run.GetField(meta.ObjectId, str)
			}
		}
	}
	return
}
