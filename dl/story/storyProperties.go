package story

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/rt"
)

type PropertySlot interface {
	GetParam() eph.EphParams
}

func (op *NamedProperty) getParam(aff string, init rt.Assignment) eph.EphParams {
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{aff},
		Class:     op.Type,
		Initially: init,
	}
}

func (op *AspectProperty) GetParam() eph.EphParams {
	// inform gives these the name "<noun> condition"
	// while tapestry relies on the name and class of the aspect to be the same.
	// we could only do that with an after the fact reduction, and with some additional mdl data.
	// ( ex. in case the same aspect is assigned twice, or twice at difference depths )
	return eph.AspectParam(op.Aspect)
}
func (op *BoolProperty) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromBool{Val: i}
	}
	return op.getParam(eph.Affinity_Bool, init)
}
func (op *NumberProperty) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromNum{Val: i}
	}
	return op.getParam(eph.Affinity_Number, init)
}
func (op *TextProperty) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromText{Val: i}
	}
	return op.getParam(eph.Affinity_Text, init)
}
func (op *NumListProperty) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromNumbers{Vals: i}
	}
	return op.getParam(eph.Affinity_NumList, init)
}
func (op *TextListProperty) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromTexts{Vals: i}
	}
	return op.getParam(eph.Affinity_TextList, init)
}
func (op *RecordProperty) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromRecord{Val: i}
	}
	return op.getParam(eph.Affinity_Record, init)
}
func (op *RecordListProperty) GetParam() eph.EphParams {
	var init rt.Assignment
	if i := op.Initially; i != nil {
		init = &core.FromRecords{Vals: i}
	}
	return op.getParam(eph.Affinity_RecordList, init)
}
