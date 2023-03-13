package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/eph"
)

type FieldDefinition interface {
	GetParam() (eph.EphParams, bool)
}

func (op *NothingField) GetParam() (nothing eph.EphParams, okay bool) {
	return
}

func (op *AspectField) GetParam() (eph.EphParams, bool) {
	// inform gives these the name "<noun> condition"
	// while tapestry relies on the name and class of the aspect to be the same.
	// we could only do that with an after the fact reduction, and with some additional mdl data.
	// ( ex. in case the same aspect is assigned twice, or twice at difference depths )
	return eph.AspectParam(op.Aspect), true
}

func (op *BoolField) GetParam() (eph.EphParams, bool) {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromBool{Value: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_Bool},
		Class:     op.Type,
		Initially: init,
	}, true
}

func (op *NumberField) GetParam() (eph.EphParams, bool) {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromNumber{Value: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_Number},
		Class:     op.Type,
		Initially: init,
	}, true
}

func (op *TextField) GetParam() (eph.EphParams, bool) {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromText{Value: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_Text},
		Class:     op.Type,
		Initially: init,
	}, true
}

func (op *RecordField) GetParam() (eph.EphParams, bool) {
	var init assign.Assignment // init here for records is usually going to be a pattern.
	if i := op.Initially; i != nil {
		init = &assign.FromRecord{Value: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_Record},
		Class:     op.Type,
		Initially: init,
	}, true
}

func (op *NumListField) GetParam() (eph.EphParams, bool) {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromNumList{Value: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_NumList},
		Class:     op.Type,
		Initially: init,
	}, true
}

func (op *TextListField) GetParam() (eph.EphParams, bool) {
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromTextList{Value: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_TextList},
		Class:     op.Type,
		Initially: init,
	}, true
}

func (op *RecordListField) GetParam() (eph.EphParams, bool) {
	// FIX: always assign the initializer, and use it to determine affinity
	// instead of the hard coded strings
	var init assign.Assignment
	if i := op.Initially; i != nil {
		init = &assign.FromRecordList{Value: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_RecordList},
		Class:     op.Type,
		Initially: init,
	}, true
}
