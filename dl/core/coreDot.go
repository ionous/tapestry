package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// Dot picks a value out of another value.
type Dot interface {
	Peek(rt.Runtime, g.Value) (g.Value, error)
	Poke(run rt.Runtime, target, value g.Value) error
}

func (op *AtField) Peek(run rt.Runtime, val g.Value) (ret g.Value, err error) {
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

func (op *AtField) Poke(run rt.Runtime, target, newValue g.Value) (err error) {
	if aff := target.Affinity(); !affine.HasFields(aff) {
		err = cmdError(op, errutil.New(aff, "doesn't have fields"))
	} else if field, e := safe.GetText(run, op.Field); e != nil {
		err = cmdError(op, e)
	} else if e := target.SetFieldByName(field.String(), newValue); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *AtIndex) Peek(run rt.Runtime, val g.Value) (ret g.Value, err error) {
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

func (op *AtIndex) Poke(run rt.Runtime, target, newValue g.Value) (err error) {
	if aff := target.Affinity(); !affine.IsList(aff) {
		err = cmdError(op, errutil.New(aff, "isn't a list"))
	} else if idx, e := safe.GetNumber(run, op.Index); e != nil {
		err = cmdError(op, e)
	} else if i, e := safe.Range(idx.Int()-1, 0, target.Len()); e != nil {
		err = cmdError(op, e)
	} else if e := target.SetIndex(i, newValue); e != nil {
		err = cmdError(op, e)
	}
	return
}

func Peek(run rt.Runtime, val g.Value, dot []Dot) (ret g.Value, err error) {
	for _, p := range dot {
		if next, e := p.Peek(run, val); e != nil {
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
