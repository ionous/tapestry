package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// PathEval picks a value out of another value.
type PathEval interface {
	PickValue(rt.Runtime, g.Value) (g.Value, error)
}

func (op *AtField) PickValue(run rt.Runtime, val g.Value) (ret g.Value, err error) {
	if aff := val.Affinity(); !affine.HasFields(aff) {
		err = cmdError(op, errutil.New(aff, "doesn't have fields"))
	} else if field, e := safe.GetText(run, op.Field); e != nil {
		err = cmdError(op, e)
	} else if el, e := val.FieldByName(field.String()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = el
	}
	return
}

func (op *AtIndex) PickValue(run rt.Runtime, val g.Value) (ret g.Value, err error) {
	if aff := val.Affinity(); !affine.IsList(aff) {
		err = cmdError(op, errutil.New(aff, "isn't a list"))
	} else if idx, e := safe.GetNumber(run, op.Index); e != nil {
		err = cmdError(op, e)
	} else if i, e := safe.Range(idx.Int()-1, 0, val.Len()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val.Index(i)
	}
	return
}

func PickValue(run rt.Runtime, val g.Value, path []PathEval) (ret g.Value, err error) {
	for _, p := range path {
		if next, e := p.PickValue(run, val); e != nil {
			err = e
			break
		} else {
			val = next
		}
	}
	if err == nil {
		ret = val
	}
	return
}
