package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type FieldDefinition interface {
	// since field creation is delayed anyway, rather than generate an error,
	// field info stores an error for later reporting
	GetFieldInfo(rt.Runtime) mdl.FieldInfo
}

func reduceFields(run rt.Runtime, fd []FieldDefinition) []mdl.FieldInfo {
	out := make([]mdl.FieldInfo, len(fd))
	for i, fd := range fd {
		if fd != nil {
			out[i] = fd.GetFieldInfo(run)
		}
	}
	return out
}

func (op *NothingField) GetFieldInfo(run rt.Runtime) (_ mdl.FieldInfo) {
	return
}

func (op *AspectField) GetFieldInfo(run rt.Runtime) (ret mdl.FieldInfo) {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &call.FromText{Value: i}
	}
	return defineField(run, op.AspectName, op.AspectName, affine.Text, init)
}

func (op *BoolField) GetFieldInfo(run rt.Runtime) mdl.FieldInfo {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &call.FromBool{Value: i}
	}
	return defineField(run, op.FieldName, nil, affine.Bool, init)
}

func (op *NumField) GetFieldInfo(run rt.Runtime) mdl.FieldInfo {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &call.FromNum{Value: i}
	}
	return defineField(run, op.FieldName, nil, affine.Num, init)
}

func (op *TextField) GetFieldInfo(run rt.Runtime) mdl.FieldInfo {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &call.FromText{Value: i}
	}
	return defineField(run, op.FieldName, op.KindName, affine.Text, init)
}

func (op *RecordField) GetFieldInfo(run rt.Runtime) mdl.FieldInfo {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &call.FromRecord{Value: i}
	}
	return defineField(run, op.FieldName, op.RecordName, affine.Record, init)
}

func (op *NumListField) GetFieldInfo(run rt.Runtime) mdl.FieldInfo {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &call.FromNumList{Value: i}
	}
	return defineField(run, op.FieldName, nil, affine.NumList, init)
}

func (op *TextListField) GetFieldInfo(run rt.Runtime) mdl.FieldInfo {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &call.FromTextList{Value: i}
	}
	return defineField(run, op.FieldName, op.KindName, affine.TextList, init)
}

func (op *RecordListField) GetFieldInfo(run rt.Runtime) mdl.FieldInfo {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &call.FromRecordList{Value: i}
	}
	return defineField(run, op.FieldName, op.RecordName, affine.RecordList, init)
}

func defineField(run rt.Runtime, name, cls rt.TextEval, aff affine.Affinity, init rt.Assignment) (ret mdl.FieldInfo) {
	if name, e := safe.GetText(run, name); e != nil {
		ret = mdl.FieldInfo{Error: e}
	} else if cls, e := safe.GetOptionalText(run, cls, ""); e != nil {
		ret = mdl.FieldInfo{Error: e}
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
