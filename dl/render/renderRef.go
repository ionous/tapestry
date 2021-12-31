package render

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *RenderRef) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getAssignedValue(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderRef) Affinity() (ret affine.Affinity) { return }

// GetText handles unpacking a text variable,
func (op *RenderRef) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getBool(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderRef) getBool(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getAssignedValue(run); e != nil {
		err = e
	} else if aff := v.Affinity(); aff == affine.Bool {
		ret = v
	} else {
		err = errutil.Fmt("unexpected %q", aff)
	}
	return
}

// GetText handles unpacking a text variable,
// turning an object variable into an id, or
// looking for an object of the passed name ( if no variable of the name exists. )
func (op *RenderRef) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getNum(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderRef) getNum(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getAssignedValue(run); e != nil {
		err = e
	} else if aff := v.Affinity(); aff == affine.Number {
		ret = v
	} else {
		err = errutil.Fmt("unexpected %q", aff)
	}
	return
}

// GetText handles unpacking a text variable,
// turning an object variable into an id, or
// looking for an object of the passed name ( if no variable of the name exists. )
func (op *RenderRef) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getText(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderRef) getText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.getAssignedValue(run); e != nil {
		err = e
	} else if aff := v.Affinity(); aff == affine.Text {
		ret = v
	} else if aff == affine.Object {
		ret = g.ObjectAsText(v)
	} else {
		err = errutil.Fmt("unexpected %q", aff)
	}
	return
}

func (op *RenderRef) getAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	flags := op.Flags.ToFlags()
	if val, e := getVariable(run, op.Name, flags); e != nil {
		err = e
	} else if val != nil {
		ret = val
	} else if !flags.tryObject() {
		err = g.UnknownVariable(op.Name.String())
	} else if obj, e := safe.ObjectFromString(run, op.Name.String()); e != nil {
		err = e
	} else {
		ret = obj
	}
	return
}

// returns nil if the named variable doesnt exist; errors only on critical errors.
func getVariable(run rt.Runtime, n core.VariableName, flags TryAsNoun) (ret g.Value, err error) {
	if flags.tryVariable() {
		ret, err = safe.CheckVariable(run, n.String(), "")
		if _, isUnknown := err.(g.Unknown); isUnknown {
			err = nil // simplify caller check.
		}
	}
	return
}
