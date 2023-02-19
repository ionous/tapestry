package story

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
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
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_Bool},
		Class:     op.Type,
		Initially: core.AssignFromBool(op.Initially),
	}
}

func (op *NumberField) GetParam() eph.EphParams {
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_Number},
		Class:     op.Type,
		Initially: core.AssignFromNumber(op.Initially),
	}
}

func (op *TextField) GetParam() eph.EphParams {
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_Text},
		Class:     op.Type,
		Initially: core.AssignFromText(op.Initially),
	}
}

func (op *NumListField) GetParam() eph.EphParams {
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_NumList},
		Class:     op.Type,
		Initially: core.AssignFromNumList(op.Initially),
	}
}

func (op *TextListField) GetParam() eph.EphParams {
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_TextList},
		Class:     op.Type,
		Initially: core.AssignFromTextList(op.Initially),
	}
}

func (op *RecordField) GetParam() eph.EphParams {
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_Record},
		Class:     op.Type,
		Initially: core.AssignFromRecord(op.Initially),
	}
}

func (op *RecordListField) GetParam() eph.EphParams {
	return eph.EphParams{
		Name:      op.Name,
		Affinity:  eph.Affinity{eph.Affinity_RecordList},
		Class:     op.Type,
		Initially: core.AssignFromRecordList(op.Initially),
	}
}
