package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *SourceValue) GetValue(run rt.Runtime) (ret g.Value, err error) {
	switch src := op.Value.(type) {
	case *FromBool:
		ret, err = safe.GetBool(run, src.Val)
	case *FromNumber:
		ret, err = safe.GetNumber(run, src.Val)
	case *FromText:
		ret, err = safe.GetText(run, src.Val)
	case *FromList:
		ret, err = safe.GetList(run, src.Val)
	case *FromRecord:
		ret, err = safe.GetRecord(run, src.Val)
	default:
		err = safe.MissingEval("source value")
	}
	return
}

func (op *SetVarFromValue) Execute(run rt.Runtime) (err error) {
	if val, e := op.Value.GetValue(run); e != nil {
		err = cmdError(op, e)
	} else if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if e := SetVariable(run, name.String(), op.Dot, val); e != nil {
		err = cmdError(op, e)
	}
	return
}

func SetVariable(run rt.Runtime, name string, dots []Dot, newValue g.Value) (err error) {
	if newValue == nil {
		err = safe.MissingEval("variable value")
	} else {
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
			if newValue == nil {
				err = errutil.New("set variable new value was nil!")
			} else {
				// affinity gets validated by SetField; fix: conversion.
				err = run.SetField(meta.Variables, name, newValue)
			}
		}
	}
	return
}
