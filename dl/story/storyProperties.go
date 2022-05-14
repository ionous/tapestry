package story

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/rt"
)

type Field interface {
	GetParam() eph.EphParams
}

func (op *AspectField) GetParam() eph.EphParams {
	// inform gives these the name "<noun> condition"
	// while tapestry relies on the name and class of the aspect to be the same.
	// we could only do that with an after the fact reduction, and with some additional mdl data.
	// ( ex. in case the same aspect is assigned twice, or twice at difference depths )
	return eph.AspectParam(op.Aspect)
}
func (op *BoolField) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromBool{Val: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_Bool},
		Class:     op.Type,
		Initially: init,
	}
}

func (op *NumberField) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromNum{Val: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_Number},
		Class:     op.Type,
		Initially: init,
	}
}

func (op *TextField) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromText{Val: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_Text},
		Class:     op.Type,
		Initially: init,
	}
}

func (op *NumListField) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromNumbers{Vals: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_NumList},
		Class:     op.Type,
		Initially: init,
	}
}

func (op *TextListField) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromTexts{Vals: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_TextList},
		Class:     op.Type,
		Initially: init,
	}
}

func (op *RecordField) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromRecord{Val: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_Record},
		Class:     op.Type,
		Initially: init,
	}
}

func (op *RecordListField) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromRecords{Vals: i}
	}
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_RecordList},
		Class:     op.Type,
		Initially: init,
	}
}
