package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type FieldDefinition interface {
	FieldInfo() mdl.FieldInfo
}

func (op *NothingField) FieldInfo() (_ mdl.FieldInfo) {
	return
}

func (op *AspectField) FieldInfo() mdl.FieldInfo {
	// inform gives these the name "<noun> condition"
	// while tapestry relies on the name and class of the aspect to be the same.
	// we could only do that with an after the fact reduction, and with some additional mdl data.
	// ( ex. in case the same aspect is assigned twice, or twice at difference depths )
	return mdl.FieldInfo{
		Name:     lang.Normalize(op.Aspect),
		Class:    lang.Normalize(op.Aspect),
		Affinity: affine.Text,
	}
}

func (op *BoolField) FieldInfo() mdl.FieldInfo {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromBool{Value: i}
	}
	return mdl.FieldInfo{
		Name:     lang.Normalize(op.Name),
		Class:    lang.Normalize(op.Type),
		Affinity: affine.Bool,
		Init:     init,
	}
}

func (op *NumberField) FieldInfo() mdl.FieldInfo {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromNumber{Value: i}
	}
	return mdl.FieldInfo{
		Name:     lang.Normalize(op.Name),
		Class:    lang.Normalize(op.Type),
		Affinity: affine.Number,
		Init:     init,
	}
}

func (op *TextField) FieldInfo() mdl.FieldInfo {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromText{Value: i}
	}
	return mdl.FieldInfo{
		Name:     lang.Normalize(op.Name),
		Class:    lang.Normalize(op.Type),
		Affinity: affine.Text,
		Init:     init,
	}

}

func (op *RecordField) FieldInfo() mdl.FieldInfo {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromRecord{Value: i}
	}
	return mdl.FieldInfo{
		Name:     lang.Normalize(op.Name),
		Class:    lang.Normalize(op.Type),
		Affinity: affine.Record,
		Init:     init,
	}
}

func (op *NumListField) FieldInfo() mdl.FieldInfo {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromNumList{Value: i}
	}
	return mdl.FieldInfo{
		Name:     lang.Normalize(op.Name),
		Class:    lang.Normalize(op.Type),
		Affinity: affine.NumList,
		Init:     init,
	}
}

func (op *TextListField) FieldInfo() mdl.FieldInfo {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromTextList{Value: i}
	}
	return mdl.FieldInfo{
		Name:     lang.Normalize(op.Name),
		Class:    lang.Normalize(op.Type),
		Affinity: affine.TextList,
		Init:     init,
	}
}

func (op *RecordListField) FieldInfo() mdl.FieldInfo {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromRecordList{Value: i}
	}
	return mdl.FieldInfo{
		Name:     lang.Normalize(op.Name),
		Class:    lang.Normalize(op.Type),
		Affinity: affine.RecordList,
		Init:     init,
	}
}
