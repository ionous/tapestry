package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
)

type PropertySlot interface {
	GetParam() eph.EphParams
}

func (op *NamedProperty) getParam(aff string) eph.EphParams {
	return eph.EphParams{
		Name:     op.Name,
		Affinity: eph.Affinity{aff},
		Class:    op.Type,
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
	return op.getParam(eph.Affinity_Bool)
}
func (op *NumberProperty) GetParam() eph.EphParams {
	return op.getParam(eph.Affinity_Number)
}
func (op *TextProperty) GetParam() eph.EphParams {
	return op.getParam(eph.Affinity_Text)
}
func (op *NumListProperty) GetParam() eph.EphParams {
	return op.getParam(eph.Affinity_NumList)
}
func (op *TextListProperty) GetParam() eph.EphParams {
	return op.getParam(eph.Affinity_TextList)
}
func (op *RecordProperty) GetParam() eph.EphParams {
	return op.getParam(eph.Affinity_Record)
}
func (op *RecordListProperty) GetParam() eph.EphParams {
	return op.getParam(eph.Affinity_RecordList)
}
