package render

import (
	"strconv"
	"strings"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

func (op *RenderName) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getName(run); e != nil {
		err = cmdError(op, e)
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
		switch v, e := run.GetField(object.Variables, op.Name); e.(type) {
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

			case affine.Object:
				ret, err = op.getPrintedNamedOf(run, v.String())

			case affine.Text:
				if n := v.String(); strings.HasPrefix(n, "#") {
					// if its an object id, get its printed name
					ret, err = op.getPrintedNamedOf(run, n)
				} else {
					// if its not, just assume the author was asking for the variable's text
					ret = v
				}
			}
		}
	}
	return
}

func (op *RenderName) getPrintedNamedOf(run rt.Runtime, objectName string) (ret g.Value, err error) {
	if printedName, e := safe.GetText(run, &core.BufferText{core.MakeActivity(
		&core.CallPattern{
			Pattern:   value.PatternName{Str: "print_name"},
			Arguments: core.Args(&core.FromText{&core.TextValue{value.Text{Str: objectName}}})})}); e != nil {
		err = e
	} else {
		ret = printedName
	}
	return
}
