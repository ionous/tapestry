package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// IntoTargetFields: part of PutAtField
type IntoTargetFields interface {
	GetTargetFields(run rt.Runtime) (g.Value, error)
}

// GetTargetFields returns an object supporting field access.
func (op *IntoObj) GetTargetFields(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if v == nil {
		err = cmdError(op, g.NothingObject)
	} else {
		ret = v
	}
	return
}

// GetTargetFields returns a record or object supporting field access.
// ( see also FromVar )
func (op *IntoVar) GetTargetFields(run rt.Runtime) (ret g.Value, err error) {
	if v, e := fieldsFromVar(run, op.Var.String()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func GetTargetFields(run rt.Runtime, fields IntoTargetFields) (ret g.Value, err error) {
	if fields == nil {
		err = safe.MissingEval("target fields")
	} else {
		ret, err = fields.GetTargetFields(run)
	}
	return
}
