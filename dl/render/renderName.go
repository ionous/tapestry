package render

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/inflect/en"
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
	if name := op.Name; en.IsCapitalized(name) {
		ret, err = op.getPrintedObjectName(run, name)

	} else {
		// first check if there's a variable of the requested name
		switch v, e := run.GetField(meta.Variables, op.Name); e.(type) {
		default:
			err = e
		case g.Unknown:
			// if there was no such variable, then it's probably an object name
			ret, err = op.getPrintedObjectName(run, name)

		case nil:
			// trying to print a variable? what kind?
			switch aff := v.Affinity(); aff {
			default:
				err = errutil.Fmt("can't render name of variable %q a %s", op.Name, aff)

			case affine.Bool:
				str := strconv.FormatBool(v.Bool())
				ret = g.StringOf(str)

			case affine.Number:
				str := strconv.FormatFloat(v.Float(), 'g', -1, 64)
				ret = g.StringOf(str)

			case affine.Text:
				// if there's no type, just assume the author was asking for the variable's text
				// if the string is empty: allow it to print nothing... backwards compat for printing nil objects
				if str, kind := v.String(), v.Type(); len(kind) == 0 || len(str) == 0 {
					ret = v
				} else if k, e := run.GetKindByName(kind); e != nil {
					err = e
				} else if b := g.Base(k); b != kindsOf.Kind.String() {
					ret = v
				} else {
					ret, err = op.getPrintedValue(run, str, kind)
				}
			}
		}
	}
	return
}

func (op *RenderName) getPrintedObjectName(run rt.Runtime, name string) (ret g.Value, err error) {
	if obj, e := run.GetField(meta.ObjectId, name); e != nil {
		err = e
	} else {
		ret, err = op.getPrintedValue(run, obj.String(), obj.Type())
	}
	return
}

func (op *RenderName) getPrintedValue(run rt.Runtime, n, k string) (ret g.Value, err error) {
	if printedName, e := safe.GetText(run, &core.BufferText{Exe: core.MakeActivity(
		&assign.CallPattern{
			PatternName: "print name",
			Arguments: core.MakeArgs(&assign.FromText{Value: &literal.TextValue{
				Value: n,
				Kind:  k,
			}})})}); e != nil {
		err = e
	} else {
		ret = printedName
	}
	return
}
