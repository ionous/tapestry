package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *SetVarFromBool) Execute(run rt.Runtime) (err error) {
	if newValue, e := safe.GetBool(run, op.Bool); e != nil {
		err = cmdError(op, e)
	} else if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if e := SetVariable(run, name.String(), op.Dot, newValue); e != nil {
		err = cmdError(op, e)
	}
	return
}
func (op *SetVarFromNumber) Execute(run rt.Runtime) (err error) {
	if newValue, e := safe.GetNumber(run, op.Number); e != nil {
		err = cmdError(op, e)
	} else if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if e := SetVariable(run, name.String(), op.Dot, newValue); e != nil {
		err = cmdError(op, e)
	}
	return
}
func (op *SetVarFromText) Execute(run rt.Runtime) (err error) {
	if newValue, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if e := SetVariable(run, name.String(), op.Dot, newValue); e != nil {
		err = cmdError(op, e)
	}
	return
}
func (op *SetVarFromList) Execute(run rt.Runtime) (err error) {
	if newValue, e := safe.GetList(run, op.List); e != nil {
		err = cmdError(op, e)
	} else if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if e := SetVariable(run, name.String(), op.Dot, newValue); e != nil {
		err = cmdError(op, e)
	}
	return
}
func (op *SetVarFromRecord) Execute(run rt.Runtime) (err error) {
	if newValue, e := safe.GetRecord(run, op.Record); e != nil {
		err = cmdError(op, e)
	} else if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if e := SetVariable(run, name.String(), op.Dot, newValue); e != nil {
		err = cmdError(op, e)
	}
	return
}

func SetVariable(run rt.Runtime, name string, dots []Dot, newValue g.Value) (err error) {
	// unpack the dotted path; fix: using pointers would be better.
	if cnt := len(dots); cnt > 0 {
		// record up to the penultimate value.
		// ie. if we want x.y.z = w; while we have to read y from x;
		// we don't need to read z from y because it's getting replaced.
		if val, e := run.GetField(meta.Variables, name); e != nil {
			err = e
		} else {
			last := cnt - 1
			vs := make([]g.Value, 0, last)
			for i := 0; i < last; i++ {
				at := dots[i]
				if next, e := at.Peek(run, val); e != nil {
					err = e
					break
				} else {
					vs[i] = next
					val = next
				}
			}
			// want: a.b.c.d = newValue; cnt=3
			// read: vs, val = { a.b, a.b.c }
			// poke: a.b.c[d] = newValue
			// poke: a.b[c]   = a.b.c
			// poke: a[b]     = a.b
			// set: a= newValue
			if err == nil {
				for {
					at := dots[last]
					if e := at.Poke(run, val, newValue); e != nil {
						err = e
						break
					} else {
						newValue = val
						if last--; last < 0 {
							break
						} else {
							val = vs[last]
						}
					}
				}
			}
		}
	}
	// finally: a= newValue
	if err == nil {
		// affinity gets validated by SetField; fix: conversion.
		err = run.SetField(meta.Variables, name, newValue)
	}
	return
}
