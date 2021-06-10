package core

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

func (op *GetAtField) Affinity() affine.Affinity { return "" }

// GetAssignedValue implements Assignment so we can SetXXX values from variables without a FromXXX statement in between.
func (op *GetAtField) GetAssignedValue(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, "")
}

func (op *GetAtField) GetBool(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Bool)
}

func (op *GetAtField) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Number)
}

func (op *GetAtField) GetText(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Text)
}

func (op *GetAtField) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.Record)
}

func (op *GetAtField) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.NumList)
}

func (op *GetAtField) GetTextList(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.TextList)
}

func (op *GetAtField) GetRecordList(run rt.Runtime) (g.Value, error) {
	return op.unpack(run, affine.RecordList)
}

func (op *GetAtField) unpack(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if src, e := GetSourceFields(run, op.From); e != nil {
		err = cmdError(op, e)
	} else if v, e := safe.Unpack(src, op.Field, aff); e != nil {
		err = errutil.Fmt("trying field %q %w", op.Field, cmdError(op, e))
	} else {
		ret = v
	}
	return
}
