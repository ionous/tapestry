package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *SetObjFromValue) Execute(run rt.Runtime) (err error) {
	if val, e := op.Value.GetValue(run); e != nil {
		err = cmdError(op, e)
	} else if name, e := safe.GetText(run, op.Name); e != nil {
		err = cmdError(op, e)
	} else if field, e := safe.GetText(run, op.Field); e != nil {
		err = cmdError(op, e)
	} else if obj, e := safe.ObjectFromString(run, name.String()); e != nil {
		err = cmdError(op, e)
	} else if e := obj.SetFieldByName(field.String(), val); e != nil {
		err = cmdError(op, e)
	}
	return
}
