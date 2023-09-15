package safe

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
)

// used when converting values to fields that might require objects.
// if the target field (ex. a pattern local) requires text of a certain type
// and the incoming value is untyped, try to convert it.
func RectifyText(run rt.Runtime, ft g.Field, val g.Value) (ret g.Value, err error) {
	ret = val // provisionally.
	// assigning to a field of typed text (which refers to an object?)
	if ft.Affinity == affine.Text && len(ft.Type) > 0 {
		// and the input is untyped?
		// (tbd: reject objects of incompatible type?)
		if val.Affinity() == affine.Text && val.Len() > 0 && len(val.Type()) == 0 {
			// is the target field of object type?
			if k, e := run.GetKindByName(ft.Type); e != nil {
				err = e
			} else if k.Implements(kindsOf.Kind.String()) {
				ret, err = run.GetField(meta.ObjectId, val.String())
			}
		}
	}
	return
}
