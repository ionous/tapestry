package core

import (
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func (op *PutAtField) Execute(run rt.Runtime) (err error) {
	if e := op.pack(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *PutAtField) pack(run rt.Runtime) (err error) {
	if val, e := safe.GetAssignedValue(run, op.From); e != nil {
		err = e
	} else if target, e := GetTargetFields(run, op.Into); e != nil {
		err = e
	} else {
		err = target.SetFieldByName(op.AtField.Value(), val)
	}
	return
}
