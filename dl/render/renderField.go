package render

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

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
		ret, err = run.GetField(meta.ObjectValue, name)
	} else {
		// try as a variable:
		switch v, e := run.GetField(meta.Variables, name); e.(type) {
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
			ret, err = run.GetField(meta.ObjectValue, name)
		default:
			err = e
		}
	}
	return
}
