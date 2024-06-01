package render

import (
	"fmt"
	"strconv"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/call"

	"git.sr.ht/~ionous/tapestry/dl/format"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

func (op *RenderName) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.getName(run); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderName) getName(run rt.Runtime) (ret rt.Value, err error) {
	// uppercase names are assumed to be requests for object names.
	if name := op.Name; inflect.IsCapitalized(name) {
		ret, err = op.getPrintedObjectName(run, name)

	} else {
		// first check if there's a variable of the requested name
		switch v, e := run.GetField(meta.Variables, op.Name); e.(type) {
		default:
			err = e
		case rt.Unknown:
			// if there was no such variable, then it's probably an object name
			ret, err = op.getPrintedObjectName(run, name)

		case nil:
			// trying to print a variable? what kind?
			switch aff := v.Affinity(); aff {
			default:
				err = fmt.Errorf("can't render name of variable %q a %s", op.Name, aff)

			case affine.Bool:
				str := strconv.FormatBool(v.Bool())
				ret = rt.StringOf(str)

			case affine.Num:
				str := strconv.FormatFloat(v.Float(), 'g', -1, 64)
				ret = rt.StringOf(str)

			case affine.Text:
				// if there's no type, just assume the author was asking for the variable's text
				// if the string is empty: allow it to print nothing... backwards compat for printing nil objects
				if str, kind := v.String(), v.Type(); len(kind) == 0 || len(str) == 0 {
					ret = v
				} else if k, e := run.GetKindByName(kind); e != nil {
					err = e
				} else if !k.Implements(kindsOf.Kind.String()) {
					ret = v
				} else {
					ret, err = op.getPrintedValue(run, str, kind)
				}
			}
		}
	}
	return
}

func (op *RenderName) getPrintedObjectName(run rt.Runtime, name string) (ret rt.Value, err error) {
	if obj, e := run.GetField(meta.ObjectId, name); e != nil {
		err = e
	} else {
		ret, err = op.getPrintedValue(run, obj.String(), obj.Type())
	}
	return
}

func (op *RenderName) getPrintedValue(run rt.Runtime, n, k string) (ret rt.Value, err error) {
	if printedName, e := safe.GetText(run, &format.BufferText{Exe: []rt.Execute{
		&call.CallPattern{
			PatternName: "print name",
			Arguments: assign.MakeArgs(&assign.FromText{Value: &literal.TextValue{
				Value: n,
				Kind:  k,
			}})}}}); e != nil {
		err = e
	} else {
		ret = printedName
	}
	return
}
