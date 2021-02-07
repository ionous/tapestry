package render

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

// RenderField -
// in template phrases, picks between record variables, object variables, and named global objects.
// ex. could be "ringBearer", "SamWise", or "frodo"
type RenderField struct {
	Name rt.TextEval
}

// Compose implements composer.Composer, although this is a
func (*RenderField) Compose() composer.Spec {
	return composer.Spec{
		Group: "internal",
	}
}

// GetSourceFields returns a value supporting field access.
func (op *RenderField) GetSourceFields(run rt.Runtime) (ret g.Value, err error) {
	if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if v, e := getSourceFields(run, name.String()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func getSourceFields(run rt.Runtime, name string) (ret g.Value, err error) {
	// uppercase names are assumed to be requests for object names.
	if lang.IsCapitalized(name) {
		ret, err = run.GetField(object.Value, name)
	} else {
		// try as a variable:
		switch v, e := run.GetField(object.Variables, name); e.(type) {
		case nil:
			// convert the variable to a set of fields
			switch aff := v.Affinity(); aff {
			case affine.Record:
				ret = v
			case affine.Text:
				ret, err = safe.ObjectFromString(run, v.String())
			default:
				err = errutil.Fmt("unexpected %q for %q", aff, name)
			}
		case g.Unknown:
			// no such variable? try as an object
			ret, err = run.GetField(object.Value, name)
		default:
			err = e
		}
	}
	return
}
