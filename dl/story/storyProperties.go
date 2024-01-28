package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type FieldDefinition interface {
	FieldInfo(rt.Runtime) (mdl.FieldInfo, error)
}

func (op *NothingField) FieldInfo(run rt.Runtime) (_ mdl.FieldInfo, _ error) {
	return
}

func (op *AspectField) FieldInfo(run rt.Runtime) (ret mdl.FieldInfo, err error) {
	return defineField(run, op.Aspect, op.Aspect, affine.Text, nil)
}

func (op *BoolField) FieldInfo(run rt.Runtime) (mdl.FieldInfo, error) {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromBool{Value: i}
	}
	return defineField(run, op.Name, op.Type, affine.Bool, init)
}

func (op *NumberField) FieldInfo(run rt.Runtime) (mdl.FieldInfo, error) {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromNumber{Value: i}
	}
	return defineField(run, op.Name, op.Type, affine.Number, init)
}

func (op *TextField) FieldInfo(run rt.Runtime) (mdl.FieldInfo, error) {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromText{Value: i}
	}
	return defineField(run, op.Name, op.Type, affine.Text, init)
}

func (op *RecordField) FieldInfo(run rt.Runtime) (mdl.FieldInfo, error) {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromRecord{Value: i}
	}
	return defineField(run, op.Name, op.Type, affine.Record, init)
}

func (op *NumListField) FieldInfo(run rt.Runtime) (mdl.FieldInfo, error) {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromNumList{Value: i}
	}
	return defineField(run, op.Name, op.Type, affine.NumList, init)
}

func (op *TextListField) FieldInfo(run rt.Runtime) (mdl.FieldInfo, error) {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromTextList{Value: i}
	}
	return defineField(run, op.Name, op.Type, affine.TextList, init)
}

func (op *RecordListField) FieldInfo(run rt.Runtime) (mdl.FieldInfo, error) {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromRecordList{Value: i}
	}
	return defineField(run, op.Name, op.Type, affine.RecordList, init)
}

func defineField(run rt.Runtime, name, cls rt.TextEval, aff affine.Affinity, init rt.Assignment) (ret mdl.FieldInfo, err error) {
	if name, e := safe.GetText(run, name); e != nil {
		err = e
	} else if cls, e := safe.GetOptionalText(run, cls, ""); e != nil {
		err = e
	} else {
		ret = mdl.FieldInfo{
			Name:     inflect.Normalize(name.String()),
			Class:    inflect.Normalize(cls.String()),
			Affinity: aff,
			Init:     init,
		}
	}
	return
}
