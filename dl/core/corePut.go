package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

/**
 * put: eval,
 * intoObj/intoObjNamed/intoRec: varName,objName,textEval,
 * atField: string.
 */
type PutAtField struct {
	From    rt.Assignment    `if:"selector"`
	Into    IntoTargetFields `if:"selector"`
	AtField string
}

func (*PutAtField) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "put", Role: composer.Command},
		Group:  "variables",
		Desc:   "Put into field: put a value into the field of an record or object",
	}
}

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
		err = target.SetFieldByName(op.AtField, val)
	}
	return
}
