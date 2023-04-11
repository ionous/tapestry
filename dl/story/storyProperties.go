package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
)

type fieldType func(name, class string, aff affine.Affinity, init assign.Assignment) error

type FieldDefinition interface {
	DeclareField(fieldType) error
}

func (op *NothingField) DeclareField(fn fieldType) (none error) {
	return
}

func (op *AspectField) DeclareField(fn fieldType) error {
	// inform gives these the name "<noun> condition"
	// while tapestry relies on the name and class of the aspect to be the same.
	// we could only do that with an after the fact reduction, and with some additional mdl data.
	// ( ex. in case the same aspect is assigned twice, or twice at difference depths )
	name, class := op.Aspect, op.Aspect
	return fn(name, class, affine.Text, nil)
}

func (op *BoolField) DeclareField(fn fieldType) error {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromBool{Value: i}
	}
	return fn(op.Name, op.Type, affine.Bool, init)
}

func (op *NumberField) DeclareField(fn fieldType) error {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromNumber{Value: i}
	}
	return fn(op.Name, op.Type, affine.Number, init)
}

func (op *TextField) DeclareField(fn fieldType) error {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromText{Value: i}
	}
	return fn(op.Name, op.Type, affine.Text, init)
}

func (op *RecordField) DeclareField(fn fieldType) error {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromRecord{Value: i}
	}
	return fn(op.Name, op.Type, affine.Record, init)
}

func (op *NumListField) DeclareField(fn fieldType) error {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromNumList{Value: i}
	}
	return fn(op.Name, op.Type, affine.NumList, init)
}

func (op *TextListField) DeclareField(fn fieldType) error {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromTextList{Value: i}
	}
	return fn(op.Name, op.Type, affine.TextList, init)
}

func (op *RecordListField) DeclareField(fn fieldType) error {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromRecordList{Value: i}
	}
	return fn(op.Name, op.Type, affine.RecordList, init)
}
