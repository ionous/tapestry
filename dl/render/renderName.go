package render

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *RenderName) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getName(run); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderName) getName(run rt.Runtime) (ret g.Value, err error) {
	// uppercase names are assumed to be requests for object names.
	if name := op.Name; lang.IsCapitalized(name) {
		ret, err = op.getPrintedNamedOf(run, name)
	} else {
		// first check if there's a variable of the requested name
		switch v, e := run.GetField(meta.Variables, op.Name); e.(type) {
		default:
			err = e
		case g.Unknown:
			// if there was no such variable, then it's probably an object name
			ret, err = op.getPrintedNamedOf(run, name)
		case nil:
			switch aff := v.Affinity(); aff {
			default:
				err = errutil.Fmt("variable %q is %s not text or object", op.Name, aff)
			case affine.Number:
				str := strconv.FormatFloat(v.Float(), 'g', -1, 64)
				ret = g.StringOf(str)

			case affine.Text:
				str := v.String()
				// if there's no type, just assume the author was asking for the variable's text
				// if the string is empty: allow it to print nothing... backwards compat for printing nil objects
				if vt := v.Type(); len(vt) == 0 || len(str) == 0 {
					ret = v
				} else if k, e := run.GetKindByName(vt); e != nil {
					err = e
				} else if k.Path()[0] != kindsOf.Kind.String() {
					ret = v
				} else {
					ret, err = op.getPrintedNamedOf(run, str)
				}
			}
		}
	}
	return
}

func (op *RenderName) getPrintedNamedOf(run rt.Runtime, objectName string) (ret g.Value, err error) {
	if printedName, e := safe.GetText(run, &core.BufferText{Does: core.MakeActivity(
		&assign.CallPattern{
			PatternName: "print_name",
			Arguments:   core.MakeArgs(&assign.FromText{Value: T(objectName)})})}); e != nil {
		err = e
	} else {
		ret = printedName
	}
	return
}
